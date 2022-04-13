// Package idleness supplies a function that intercepts gRPC requests and
// shuts down a server if a specified period elapses between requests.
package idleness

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// IdlenessMonitor monitors requests to a gRPC server and shuts down the server
// if it is idle for a specified period.
type IdlenessMonitor struct {
	// idlePeriod specifies a duration; after this amount of time passes and no
	// request is received or being processed, the server is shut down.
	idlePeriod     time.Duration
	lastActiveTime *lockedValue[time.Time]

	interceptor grpc.UnaryServerInterceptor

	shutdownServerCallback func()

	// done is closed when the monitor is disabled.
	done      chan struct{}
	closeDone func()
}

func NewIdlenessMonitor(idlePeriod time.Duration) *IdlenessMonitor {
	done := make(chan struct{})
	closeDone := onceCloseFn(done)
	timer := time.NewTimer(idlePeriod)

	type status struct {
		lastActiveTime time.Time
		enabled        bool
	}

	lastActiveTime := &lockedValue[time.Time]{
		value: time.Now(),
	}
	expireCallback := func() {}
	go func() {
		for {
			select {
			case <-done:
				return
			case <-timer.C:
				var t time.Time = lastActiveTime.read()
				if t.IsZero() {
					break
				}
				expireTime := t.Add(idlePeriod)
				timeUntilExpiration := time.Until(expireTime)
				if timeUntilExpiration <= 0 {
					expireCallback()
					break
				}
				timer.Reset(timeUntilExpiration)
			}
		}
	}()
	disableIdlenessTracking := func() {
		lastActiveTime.update(func(orig time.Time) time.Time {
			return time.Time{}
		})
	}
	interceptor := grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		disableIdlenessTracking()
		defer func() {
			lastActiveTime.update(func(orig time.Time) time.Time { return time.Now() })
			timer.Reset(idlePeriod)
		}()

		return handler(ctx, req)
	})
	monitor := &IdlenessMonitor{
		done:           done,
		closeDone:      closeDone,
		idlePeriod:     idlePeriod,
		interceptor:    interceptor,
		lastActiveTime: lastActiveTime,
	}
	expireCallback = func() {
		if monitor.shutdownServerCallback != nil {
			monitor.shutdownServerCallback()
		}
	}
	return monitor
}

func (im *IdlenessMonitor) Close() {
	im.closeDone()
}

// UnaryInterceptor returns an interceptor that can be passed as a
// grpc.ServerOption to grpc.NewServer that will keep track of
func (im *IdlenessMonitor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return im.interceptor
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
