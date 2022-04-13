// Package idleness supplies a function that intercepts gRPC requests and
// shuts down a server if a specified period elapses between requests.
package idleness

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Tracker monitors requests to a gRPC server and shuts down the server
// if it is idle for a specified period.
type Tracker struct {
	// idlePeriod specifies a duration; after this amount of time passes and no
	// request is received or being processed, the server is shut down.
	idlePeriod time.Duration
	expiration *lockedValue[expirationState]

	// The ticks fired by timer are the only times when the tracker checks to
	// see if the idle period has elapsed. See the constructor.
	timer *time.Timer

	// done is closed when the monitor is disabled.
	done      chan struct{}
	closeDone func()
}

func (im *Tracker) Close() {
	im.closeDone()
}

// UnaryInterceptor returns an interceptor that can be passed as a
// grpc.ServerOption to grpc.NewServer that will keep track of
func (im *Tracker) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		markActivityEnded := im.MarkContinuousActivity()
		defer markActivityEnded()
		return handler(ctx, req)
	}
}

// MarkActive sets the last active time of the timer to time.Now()
func (im *Tracker) MarkActive() {
	im.expiration.update(func(orig expirationState) expirationState {
		orig.lastActiveTime = time.Now()
		return orig
	})
}

// MarkContinuousActivity tracks a period of continuous activity. When the
// activity is complete, the caller should call the returned function.
//
// The expiration function will not be called while a continuous activity is
// ongoing.
func (im *Tracker) MarkContinuousActivity() (markActivityEnded func()) {
	im.expiration.update(func(orig expirationState) expirationState {
		orig.lastActiveTime = time.Now()
		orig.ongoingActivityCount++
		return orig
	})
	return func() {
		im.expiration.update(func(orig expirationState) expirationState {
			orig.lastActiveTime = time.Now()
			orig.ongoingActivityCount--
			return orig
		})
	}
}

// SetEnabled sets whether the tracker should call the expiration callback when
// the idle period has elapsed.
func (im *Tracker) SetEnabled(trackingEnabled bool) {
	changedToEnabled := false
	im.expiration.update(func(orig expirationState) expirationState {
		changedToEnabled = trackingEnabled && orig.disabled
		orig.disabled = !trackingEnabled
		return orig
	})
	if changedToEnabled {
		// Immediately check to see if the idleness period has elapsed.
		im.timer.Reset(0)
	}
}

type expirationState struct {
	lastActiveTime       time.Time
	disabled             bool
	ongoingActivityCount int
}

func NewIdlenessMonitor(idlePeriod time.Duration, idlePeriodElapsed func()) *Tracker {
	done := make(chan struct{})
	timer := time.NewTimer(idlePeriod)

	expiration := &lockedValue[expirationState]{
		value: expirationState{
			disabled:             false,
			lastActiveTime:       time.Now(),
			ongoingActivityCount: 0,
		},
	}
	go func() {
		for {
			select {
			case <-done:
				return
			// The timer channel should receive a tick at least every idlenessPeriod.
			case <-timer.C:
				state := expiration.read()
				if state.disabled {
					// No need to do anything. When the timer is enabled again,
					// it will be reset by SetEnabled().
					break // out of select
				}
				if state.ongoingActivityCount > 0 {
					timer.Reset(idlePeriod)
					break
				}
				expireTime := state.lastActiveTime.Add(idlePeriod)
				timeUntilExpiration := time.Until(expireTime)
				if timeUntilExpiration <= 0 {
					idlePeriodElapsed()
					break
				}
				timer.Reset(timeUntilExpiration)
			}
		}
	}()
	monitor := &Tracker{
		done:       done,
		closeDone:  onceCloseFn(done),
		idlePeriod: idlePeriod,
		timer:      timer,
		expiration: expiration,
	}
	return monitor
}

// onceCloseFn returns a function that closes ch once; subsequent calls to close
// do nothing.
func onceCloseFn[T any](ch chan T) func() {
	once := &sync.Once{}
	return func() {
		once.Do(func() {
			close(ch)
		})
	}
}

type lockedValue[T any] struct {
	lock  sync.RWMutex
	value T
}

func (v *lockedValue[T]) update(updater func(orig T) T) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.value = updater(v.value)
}

func (v *lockedValue[T]) read() T {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return v.value
}
