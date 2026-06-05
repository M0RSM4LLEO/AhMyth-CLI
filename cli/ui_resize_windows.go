//go:build windows
// +build windows

package cli

import "time"

// StartResizeListener starts Windows resize detection.
func StartResizeListener() {
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

// StopResizeListener cleans up the listener.
func StopResizeListener() {
	closeOnce.Do(func() {
		close(pollDone)
	})
}
