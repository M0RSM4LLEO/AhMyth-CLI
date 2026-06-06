package cli

import (
	"bufio"
	"fmt"
	"strings"
)

func drawPayloadOptions(prompt string) {
	ClearScreen()
	HideCursor()
	defer ShowCursor()

	RenderBuffered(func(w *bufio.Writer) {
		PrintLogo()
		fmt.Fprintln(w)

		title := "Payload Options"
		titleLine := fmt.Sprintf("%s%s%s%s", Bold, Blue, title, Reset)

		lineChar := "-"
		if SupportsUnicodeGlyphs() {
			lineChar = "─"
		}

		underline := strings.Repeat(lineChar, VisibleLength(title))

		padding := 0
		if VisibleLength(title) < GetWidth() {
			padding = (GetWidth() - VisibleLength(title)) / 2
		}

		pad := strings.Repeat(" ", padding)

		fmt.Fprintln(w, pad+titleLine)
		fmt.Fprintln(w, pad+underline)
		fmt.Fprintln(w)

		fmt.Fprintln(w, "1) APK")
		fmt.Fprintln(w, "0) Back")
		fmt.Fprintln(w)

		if prompt != "" {
			fmt.Fprint(w, prompt)
		}
	})
}

func showPayloadOptions(reader *bufio.Reader) {
	for {
		currentScreen = screenPayloadOptions
		currentPrompt = "Enter choice: "

		redrawCurrentScreen()

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "0":
			currentPrompt = ""
			currentScreen = screenMenu
			return

		case "1":
			showActionScreen(
				"APK",
				"This feature is not implemented yet.",
				reader,
			)

		default:
			currentPrompt = ""
			ClearScreen()
			fmt.Println("Invalid option!")
			Delay(800)
		}
	}
}
