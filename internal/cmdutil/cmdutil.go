package cmdutil

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
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
	ctx, cancel := context.WithCancel(parentCtx)
	ch := make(chan os.Signal, 1)
	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()
	if len(signals) > 0 {
		signal.Notify(ch, signals...)
	}

	return ctx, cancel
}
