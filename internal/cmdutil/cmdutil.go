package cmdutil

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

// LookupEnvOrDefault returns os.LookupEnv(key) or defaultValue if the key is
// not in the environment.
func LookupEnvOrDefault(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

// ConcatStringsArrays flattens each []string in the arguments into a string.
func ConcatStringsArrays(values ...[]string) []string {
	ret := []string{}
	for _, v := range values {
		ret = append(ret, v...)
	}
	return ret
}

// EnsureAbsolutePathFromWorkspace returns path if path is absolute; otherwise
// it is interpretted as relative to workspaceRoot, and
// filepath.Join(workspaceRoot, path) is returned.
func EnsureAbsolutePathFromWorkspace(workspaceRoot, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(workspaceRoot, path)
}

// SignalCancelledContext returns a child context and cancellation function. If
// an incoming signal occurs that is in the signals list, the cancellation
// function is called.
func SignalCancelledContext(parentCtx context.Context, signals ...os.Signal) (ctx context.Context, stop context.CancelFunc) {
	ctx, cancel := withCustomCancelErr(parentCtx)
	ch := make(chan os.Signal, 1)
	go func() {
		select {
		case sig := <-ch:
			cancel(fmt.Errorf("received signal %s (%d)", sig, sig))
		case <-ctx.Done():
		}
	}()
	if len(signals) > 0 {
		signal.Notify(ch, signals...)
	}

	return ctx, func() { cancel(nil) }
}

type overridableErrContext struct {
	err        error
	underlying context.Context
}

func (ec *overridableErrContext) Deadline() (deadline time.Time, ok bool) {
	return ec.underlying.Deadline()
}

func (ec *overridableErrContext) Done() <-chan struct{} {
	return ec.underlying.Done()
}

func (ec *overridableErrContext) Err() error {
	overriddenErr := ec.underlying.Err()
	if overriddenErr != nil && ec.err != nil {
		return ec.err
	}
	return overriddenErr
}

func (ec *overridableErrContext) Value(key any) any {
	return ec.underlying.Value(key)
}

func withCustomCancelErr(parent context.Context) (ctx context.Context, cancelWithErr func(error)) {
	underlyingCtx, simpleCancel := context.WithCancel(parent)
	outputCtx := &overridableErrContext{underlying: underlyingCtx}
	return outputCtx, func(err error) {
		outputCtx.err = err
		simpleCancel()
	}
}
