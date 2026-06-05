package cli

import (
	"bufio"
	"fmt"
	"time"
)

func ShowSplash() {
	HideCursor()
	defer ShowCursor()

	// Fast, in-place animation for ANSI-capable terminals.
	if SupportsANSIOutput() {
		ClearScreen()

		frames := []string{"│", "┃", "█", "┃", "│"}
		duration := 2500 * time.Millisecond
		frameDelay := 100 * time.Millisecond
		start := time.Now()

		render := func(f string) {
			RenderBuffered(func(w *bufio.Writer) {
				fmt.Fprint(w, "\033[H")
				fmt.Fprintln(w)
				fmt.Fprintln(w)
				fmt.Fprintln(w)
				fmt.Fprintln(w, Center(fmt.Sprintf("%s%s%s", Cyan, f, Reset)))
				fmt.Fprintln(w)
				fmt.Fprintln(w)
			})
		}

		for time.Since(start) < duration {
			for _, f := range frames {
				if time.Since(start) >= duration {
					break
				}
				render(f)
				Delay(int(frameDelay / time.Millisecond))
			}
		}

		ClearScreen()
		RenderBuffered(func(w *bufio.Writer) {
			fmt.Fprintln(w)
			fmt.Fprintln(w)
			fmt.Fprintln(w)
			fmt.Fprintln(w, Center(fmt.Sprintf("%s%sAhMyth%s", Bold, Blue, Reset)))
			fmt.Fprintln(w, Center("Android Remote Administration Tool"))
		})

		Delay(2500)
		ClearScreen()
		return
	}

	// Plain fallback.
	ClearScreen()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(Center("AhMyth"))
	fmt.Println(Center("Android Remote Administration Tool"))
	Delay(1200)
	ClearScreen()
}