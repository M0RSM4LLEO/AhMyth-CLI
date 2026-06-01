package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
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
	// Default fallback if terminal size can't be detected
	return 80
}

// simple centering helper
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

func waitForEnter(reader *bufio.Reader) {
	fmt.Print("\nPress Enter to return to the menu...")
	_, _ = reader.ReadString('\n')
}

func ShowSplash() {
	width := getWidth()

	fmt.Print(hide)
	defer fmt.Print(show)

	frames := []string{"│", "┃", "█", "┃", "│"}

	for i := 0; i < 3; i++ {
		for _, f := range frames {
			fmt.Print(clear, home)

			line := center(fmt.Sprintf("%s%s%s", cyan, f, reset), width)
			fmt.Println("\n\n\n" + line)

			time.Sleep(80 * time.Millisecond)
		}
	}

	fmt.Print(clear, home)

	title := center(fmt.Sprintf("%s%sAhMyth%s", bold, blue, reset), width)
	sub := center("Android Remote Administration Tool", width)

	fmt.Println("\n\n\n" + title)
	fmt.Println(sub)

	time.Sleep(1500 * time.Millisecond)
	fmt.Print(clear, home)
}

func drawMenu() {
	width := getWidth()

	logo := []string{
		"    _    _     __  __       _   _     ",
		"   / \\  | |__ |  \\/  |_   _| |_| |__  ",
		"  / _ \\ | '_ \\| |\\/| | | | | __| '_ \\ ",
		" / ___ \\| | | | |  | | |_| | |_| | | |",
		"/_/   \\_\\_| |_|_|  |_|\\__, |\\__|_| |_|",
		"                       |___/           ",
		"======================================",
	}

	clearScreen()

	fmt.Print(hide)
	defer fmt.Print(show)

	for _, line := range logo {
		fmt.Println(center(fmt.Sprintf("%s%s%s", bold, blue, line), width))
	}

	fmt.Println()

	menuTitle := "Main Menu"

	// Blue + bold title, then reset immediately
	fmt.Println(fmt.Sprintf("%s%s%s%s", bold, blue, menuTitle, reset))

	// Blue underline (explicit color reset after it)
	fmt.Println(fmt.Sprintf("%s%s%s", blue, strings.Repeat("—", len(menuTitle)), reset))

	fmt.Println()

	fmt.Println("1) Listener (TO BE DONE)")
	fmt.Println("2) Payload Options (TO BE DONE)")
	fmt.Println("0) Exit")
	fmt.Println()
}

func menuLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		drawMenu()

		fmt.Print("\nSelect an option: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			clearScreen()
			fmt.Println("Build selected.")
			// Placeholder for your build logic
			waitForEnter(reader)

		case "0":
			clearScreen()
			return

		default:
			fmt.Println("Invalid option.")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ShowSplash()
	menuLoop()
}

