package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type screenType int

const (
	screenMenu screenType = iota
	screenAction
	screenPayloadOptions
)

var (
	currentScreen screenType = screenMenu

	currentPrompt string

	currentActionTitle  string
	currentActionMessage string
)

func Run() {
	ShowSplash()

	StartResizeListener()
	SetRedrawFunc(redrawCurrentScreen)

	menuLoop()
}

func redrawCurrentScreen() {
	switch currentScreen {
	case screenMenu:
		drawMenu(currentPrompt)

	case screenAction:
		drawActionScreen(currentActionTitle, currentActionMessage, currentPrompt)

	case screenPayloadOptions:
		drawPayloadOptions(currentPrompt)
	}
}

func drawMenu(prompt string) {
	ClearScreen()
	HideCursor()
	defer ShowCursor()

	RenderBuffered(func(w *bufio.Writer) {
		PrintLogo()
		PrintHeader("\nMain Menu")

		fmt.Fprintln(w)
		fmt.Fprintln(w, "1) Listener (TO BE DONE)")
		fmt.Fprintln(w, "2) Payload Options")
		fmt.Fprintln(w, "0) Exit")
		fmt.Fprintln(w)

		if prompt != "" {
			fmt.Fprint(w, prompt)
		}
	})
}

func drawActionScreen(title string, message string, prompt string) {
	ClearScreen()
	HideCursor()
	defer ShowCursor()

	RenderBuffered(func(w *bufio.Writer) {
		width := GetWidth()

		titleLine := fmt.Sprintf("%s%s%s%s", Bold, Blue, title, Reset)

		lineChar := "-"
		if SupportsUnicodeGlyphs() {
			lineChar = "─"
		}

		underline := strings.Repeat(lineChar, VisibleLength(title))

		padding := 0
		if VisibleLength(title) < width {
			padding = (width - VisibleLength(title)) / 2
		}

		pad := strings.Repeat(" ", padding)

		fmt.Fprintln(w)
		fmt.Fprintln(w, pad+titleLine)
		fmt.Fprintln(w, pad+underline)
		fmt.Fprintln(w)
		fmt.Fprintln(w, message)
		fmt.Fprintln(w)

		if prompt != "" {
			fmt.Fprint(w, prompt)
		}
	})
}

func showActionScreen(title string, message string, reader *bufio.Reader) {
	currentScreen = screenAction
	currentActionTitle = title
	currentActionMessage = message
	currentPrompt = "Press Enter to return to the menu..."

	redrawCurrentScreen()

	waitForEnter(reader)

	currentPrompt = ""
	currentScreen = screenMenu
}

func waitForEnter(reader *bufio.Reader) {
	_, _ = reader.ReadString('\n')
}

func menuLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		currentScreen = screenMenu
		currentPrompt = "Enter choice: "

		redrawCurrentScreen()

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "0":
			ClearScreen()
			StopResizeListener()
			return

		case "1":
			showActionScreen(
				"Listener",
				"This feature is not implemented yet.",
				reader,
			)

		case "2":
			showPayloadOptions(reader)

		default:
			currentPrompt = ""
			ClearScreen()
			fmt.Println("Invalid option!")
			Delay(800)
		}
	}
}
