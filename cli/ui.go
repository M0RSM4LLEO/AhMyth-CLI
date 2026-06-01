package cli

import (
	"fmt"
	"strings"
	"os"
	"golang.org/x/term"
)

const (
	reset = "\033[0m"
	cyan  = "\033[36m"
	blue  = "\033[34m"
	bold  = "\033[1m"
	clear = "\033[2J"
	home  = "\033[H"
	hide  = "\033[?25l"
	show  = "\033[?25h"
)

func getWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}

	if width < 40 {
		return 40
	}

	return width
}

func center(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}

func clearScreen() {
	fmt.Print(clear, home)
}
