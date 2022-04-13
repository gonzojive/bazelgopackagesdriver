package idleness

import (
	"testing"
	"time"
)

func TestNewIdlenessMonitorSimple(t *testing.T) {
	start := time.Now()
	var finishedCalledTime time.Time
	finished := make(chan struct{})
	monitor := NewIdlenessMonitor(time.Millisecond*200, func() {
		finishedCalledTime = time.Now()
		close(finished)
	})

	<-finished
	if elapsed := finishedCalledTime.Sub(start); elapsed < time.Millisecond*200 {
		t.Errorf("finished, but elapsed = %s, want > 200ms", elapsed)
	}

	defer monitor.Close()
}

func TestNewIdlenessMonitorContinuousCalls(t *testing.T) {
	start := time.Now()
	var finishedCalledTime time.Time
	finished := make(chan struct{})
	monitor := NewIdlenessMonitor(time.Millisecond*50, func() {
		finishedCalledTime = time.Now()
		close(finished)
	})

	for i := 0; i < 5; i++ {
		go func() {
			finished := monitor.MarkContinuousActivity()
			defer finished()
			time.Sleep(time.Millisecond * 300)
		}()
	}

	<-finished
	if elapsed := finishedCalledTime.Sub(start); elapsed < time.Millisecond*350 {
		t.Errorf("finished, but elapsed = %s, want > 350ms", elapsed)
	}

	defer monitor.Close()
}
