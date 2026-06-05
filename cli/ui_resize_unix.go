//go:build !windows
// +build !windows

package cli

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

var resizeChan = make(chan os.Signal, 1)

// StartResizeListener starts Unix resize detection.
func StartResizeListener() {
	signal.Notify(resizeChan, syscall.SIGWINCH)

	go func() {
		for range resizeChan {
			handleResize()
		}
	}()

	go func() {
		ticker := time.NewTicker(600 * time.Millisecond)
		defer ticker.Stop()

		lastWidth = GetWidth()
		lastHeight = GetHeight()

		for {
			select {
			case <-pollDone:
				return
			case <-ticker.C:
				handleResize()
			}
		}
	}()
}

// StopResizeListener cleans up the listeners.
func StopResizeListener() {
	closeOnce.Do(func() {
		close(pollDone)
		signal.Stop(resizeChan)
	})
}
