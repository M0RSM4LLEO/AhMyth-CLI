package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func Run() {
	ShowSplash()
	menuLoop()
}

func waitForEnter(reader *bufio.Reader) {
	fmt.Print("\nPress Enter to return to the menu...")
	_, _ = reader.ReadString('\n')
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
		"--------------------------------------",
		"ANDROID  REMOTE  ANDMINISTRATION  TOOL",
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
	fmt.Println(fmt.Sprintf("%s%s%s%s", bold, blue, menuTitle, reset))
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
			fmt.Println("Listener selected.")
			waitForEnter(reader)

		case "2":
			clearScreen()
			fmt.Println("Payload Options selected.")
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
