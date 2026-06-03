package cli

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
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

	resizeMu      sync.Mutex
	currentRedraw func()
	lastWidth     int
	lastHeight    int
	pollDone      = make(chan struct{})
	closeOnce     sync.Once
)

// PrintSeparator prints a full-width separator line
func PrintSeparator() {
	width := GetWidth()
	if width < 1 {
		width = 80
	}

	sep := "─"

	// Optional fallback for legacy Windows consoles
	if runtime.GOOS == "windows" {
		if os.Getenv("WT_SESSION") == "" &&
			os.Getenv("ANSICON") == "" &&
			os.Getenv("ConEmuANSI") != "ON" {
			sep = "-"
		}
	}

	fmt.Println(strings.Repeat(sep, width))
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

// Shared resize handler
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
