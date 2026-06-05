//go:build windows
// +build windows

package cli

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	enableProcessedOutput           = 0x0001
	enableVirtualTerminalProcessing = 0x0004
)

var (
	kernel32                    = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode          = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode          = kernel32.NewProc("SetConsoleMode")
)

func init() {
	if enableWindowsVTMode() {
		useANSI = true

		Reset = "\033[0m"
		Red = "\033[31m"
		Green = "\033[32m"
		Blue = "\033[34m"
		Cyan = "\033[36m"
		Bold = "\033[1m"

		Clear = "\033[2J"
		Home = "\033[H"
		Hide = "\033[?25l"
		Show = "\033[?25h"
	}
}

func enableWindowsVTMode() bool {
	handles := []*os.File{os.Stdout, os.Stderr}

	for _, f := range handles {
		if f == nil {
			continue
		}

		h := syscall.Handle(f.Fd())

		var mode uint32
		r1, _, _ := procGetConsoleMode.Call(
			uintptr(h),
			uintptr(unsafe.Pointer(&mode)),
		)
		if r1 == 0 {
			continue
		}

		mode |= enableProcessedOutput | enableVirtualTerminalProcessing

		r1, _, _ = procSetConsoleMode.Call(
			uintptr(h),
			uintptr(mode),
		)
		if r1 != 0 {
			return true
		}
	}

	return false
}