package exec

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Loop runs the runs a function f periodically
func Loop(every time.Duration, f func()) {
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(every)
	defer ticker.Stop()

	f() // run *now*

	for {
		select {
		case <-shutdownSignal:
			return
		case <-ticker.C:
			f()
		}
	}
}
