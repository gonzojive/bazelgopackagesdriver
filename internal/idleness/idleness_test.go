// Package grpckeepalive supplies a function that intercepts gRPC requests and
// shuts down a server if a specified period elapses between requests.
package idleness

import (
	"testing"
	"time"
)

func TestNewIdlenessMonitor(t *testing.T) {
	monitor := NewIdlenessMonitor(time.Millisecond * 700)
	defer monitor.Close()

}
