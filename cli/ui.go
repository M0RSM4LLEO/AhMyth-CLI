package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

var (
	ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

	resizeMu      sync.Mutex
	currentRedraw func()
	lastWidth     int
	lastHeight    int
	pollDone      = make(chan struct{})
	closeOnce     sync.Once

	useANSI    = supportsANSI()
	useUnicode = supportsUnicodeBoxDrawing()

	Reset = ""
	Red   = ""
	Green = ""
	Blue  = ""
	Cyan  = ""
	Bold  = ""

	Clear = ""
	Home  = ""
	Hide  = ""
	Show  = ""
)

func init() {
	if useANSI {
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

func SupportsANSIOutput() bool {
	return useANSI
}

func SupportsUnicodeGlyphs() bool {
	return useUnicode
}

func supportsANSI() bool {
	if runtime.GOOS != "windows" {
		return true
	}

	if os.Getenv("WT_SESSION") != "" ||
		os.Getenv("ANSICON") != "" ||
		strings.EqualFold(os.Getenv("ConEmuANSI"), "ON") {
		return true
	}

	termEnv := strings.ToLower(os.Getenv("TERM"))
	if strings.Contains(termEnv, "xterm") ||
		strings.Contains(termEnv, "ansi") ||
		strings.Contains(termEnv, "cygwin") ||
		strings.Contains(termEnv, "msys") ||
		strings.Contains(termEnv, "vt100") {
		return true
	}

	return false
}

func supportsUnicodeBoxDrawing() bool {
	if runtime.GOOS != "windows" {
		return true
	}

	if os.Getenv("WT_SESSION") != "" || os.Getenv("ANSICON") != "" {
		return true
	}

	termEnv := strings.ToLower(os.Getenv("TERM"))
	if strings.Contains(termEnv, "xterm") ||
		strings.Contains(termEnv, "ansi") ||
		strings.Contains(termEnv, "cygwin") ||
		strings.Contains(termEnv, "msys") {
		return true
	}

	return false
}

func terminalSize() (int, int) {
	fds := []int{
		int(os.Stdin.Fd()),
		int(os.Stdout.Fd()),
		int(os.Stderr.Fd()),
	}

	for _, fd := range fds {
		w, h, err := term.GetSize(fd)
		if err == nil && w > 0 && h > 0 {
			return w, h
		}
	}

	if runtime.GOOS == "windows" {
		cols, err1 := strconv.Atoi(strings.TrimSpace(os.Getenv("COLUMNS")))
		lines, err2 := strconv.Atoi(strings.TrimSpace(os.Getenv("LINES")))
		if err1 == nil && err2 == nil && cols > 0 && lines > 0 {
			return cols, lines
		}
	}

	return 80, 24
}

// PrintSeparator prints a full-width separator line.
func PrintSeparator() {
	width := GetWidth()
	if width < 1 {
		width = 80
	}

	sep := "-"
	if useUnicode {
		sep = "─"
	}

	fmt.Println(strings.Repeat(sep, width))
}

// GetWidth returns current terminal width.
func GetWidth() int {
	w, _ := terminalSize()
	return w
}

// GetHeight returns current terminal height.
func GetHeight() int {
	_, h := terminalSize()
	return h
}

// SetRedrawFunc registers the function to call on resize.
func SetRedrawFunc(f func()) {
	resizeMu.Lock()
	currentRedraw = f
	resizeMu.Unlock()
}

// handleResize is shared by the OS-specific resize listener files.
func handleResize() {
	resizeMu.Lock()

	currentW, currentH := terminalSize()
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

func RenderBuffered(render func(w *bufio.Writer)) {
	w := bufio.NewWriterSize(os.Stdout, 8192)
	render(w)
	_ = w.Flush()
}

func ClearScreen() {
	if useANSI {
		fmt.Print("\033[2J\033[3J\033[H")
		return
	}

	fmt.Print(strings.Repeat("\n", GetHeight()))
}

func Delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func HideCursor() {
	if useANSI {
		fmt.Print(Hide)
	}
}

func ShowCursor() {
	if useANSI {
		fmt.Print(Show)
	}
}

func PrintHeader(text string) {
	visibleLen := VisibleLength(text)
	if visibleLen < 1 {
		visibleLen = len(text)
	}

	lineChar := "-"
	if useUnicode {
		lineChar = "─"
	}

	RenderBuffered(func(w *bufio.Writer) {
		fmt.Fprintln(w, text)
		fmt.Fprintln(w, strings.Repeat(lineChar, visibleLen))
	})
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
		"ANDROID  REMOTE  ADMINISTRATION  TOOL ",
		"           CLI  v1.0-beta.0           ",
		"======================================",
	}

	RenderBuffered(func(w *bufio.Writer) {
		if useANSI {
			for _, line := range logo {
				fmt.Fprintln(w, Center(fmt.Sprintf("%s%s%s%s", Bold, Blue, line, Reset)))
			}
			return
		}

		for _, line := range logo {
			fmt.Fprintln(w, Center(line))
		}
	})
}
