package cli  
  
import (  
	"fmt"  
	"os"  
	"regexp"  
	"strings"  
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
  
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)  
  
// GetWidth returns current terminal width with fallback  
func GetWidth() int {  
	width, _, err := term.GetSize(int(os.Stdout.Fd()))  
	if err != nil || width < 40 {  
		return 80  
	}  
	return width  
}  
  
// VisibleLength returns length without ANSI escape codes  
func VisibleLength(s string) int {  
	return len(ansiRegex.ReplaceAllString(s, ""))  
}  
  
// Center text using current terminal width  
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
  
// High-level UI helpers  
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
