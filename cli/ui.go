package cli

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[34m"
	Cyan  = "\033[36m"
	Bold  = "\033[1m"

	Clear = "\033[2J"
	Home  = "\033[H"
	Hide  = "\033[?25l"
	Show  = "\033[?25h"
)

var (
	ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

	stopOnce      sync.Once
	resizeMu      sync.Mutex
	currentRedraw func()
	lastWidth     int
	lastHeight    int
	resizeChan    = make(chan os.Signal, 1)
	pollDone      = make(chan struct{})
)

// PrintSeparator prints a full-width separator line
func PrintSeparator() {
	width := GetWidth()
	fmt.Println(strings.Repeat("─", width))
}

// GetWidth returns current terminal width
func GetWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w < 40 {
		return 80
	}
	return w
}

// GetHeight returns current terminal height
func GetHeight() int {
	_, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || h < 10 {
		return 24
	}
	return h
}

// SetRedrawFunc registers the function to call on resize
func SetRedrawFunc(f func()) {
	resizeMu.Lock()
	currentRedraw = f
	resizeMu.Unlock()
}

// StartResizeListener starts cross-platform resize detection
func StartResizeListener() {
	// Unix-like (Linux, macOS, Termux)
	if runtime.GOOS != "windows" {
		signal.Notify(resizeChan, syscall.SIGWINCH)
		go func() {
			for range resizeChan {
				handleResize()
			}
		}()
	}

	// Polling fallback (works on Windows + extra safety)
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

// handleResize checks if size changed and triggers redraw
func handleResize() {
	resizeMu.Lock()

	currentW := GetWidth()
	currentH := GetHeight()

	changed := currentW != lastWidth || currentH != lastHeight
	if changed {
		lastWidth = currentW
		lastHeight = currentH
	}

	redraw := currentRedraw
	resizeMu.Unlock()

	if changed && redraw != nil {
		time.Sleep(30 * time.Millisecond)
		redraw()
	}
}

// StopResizeListener cleans up (call on exit if needed)
func StopResizeListener() {
	stopOnce.Do(func() {
		close(pollDone)
		signal.Stop(resizeChan)
	})
}

func VisibleLength(s string) int {
	return len(ansiRegex.ReplaceAllString(s, ""))
}

func Center(text string) string {
	return CenterWithWidth(text, GetWidth())
}

func CenterWithWidth(text string, width int) string {
	visible := VisibleLength(text)
	if visible >= width {
		return text
	}
	padding := (width - visible) / 2
	return strings.Repeat(" ", padding) + text
}

func ClearScreen() {
	fmt.Print(Clear, Home)
}

func Delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func HideCursor() {
	fmt.Print(Hide)
}

func ShowCursor() {
	fmt.Print(Show)
}

func PrintHeader(text string) {
	visibleLen := VisibleLength(text)
	line := strings.Repeat("─", visibleLen)
	fmt.Println(text)
	fmt.Println(line)
}

func PrintCentered(text string) {
	fmt.Println(Center(text))
}

func PrintLogo() {
	logo := []string{
		"    _    _     __  __       _   _     ",
		"   / \\  | |__ |  \\/  |_   _| |_| |__  ",
		"  / _ \\ | '_ \\| |\\/| | | | | __| '_ \\ ",
		" / ___ \\| | | | |  | | |_| | |_| | | |",
		"/_/   \\_\\_| |_|_|  |_|\\__, |\\__|_| |_|",
		"                       |___/           ",
		"--------------------------------------",
		"ANDROID  REMOTE  ADMINISTRATION  TOOL",
		"======================================",
	}

	for _, line := range logo {
		fmt.Println(Center(fmt.Sprintf("%s%s%s", Bold, Blue, line)))
	}
}
