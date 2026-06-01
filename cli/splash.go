package cli

import (
	"fmt"
	"time"
)

func uiDelay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func ShowSplash() {
	width := getWidth()

	fmt.Print(hide)
	defer fmt.Print(show)

	frames := []string{"│", "┃", "█", "┃", "│"}

	render := func(f string) {
		clearScreen()
		line := center(fmt.Sprintf("%s%s%s", cyan, f, reset), width)
		fmt.Println("\n\n\n" + line)
	}

	// FLASH BY TIME (MAIN CHANGE)
	duration := 3000 * time.Millisecond // total flash time
	frameDelay := 100 * time.Millisecond

	start := time.Now()

	for time.Since(start) < duration {
		for _, f := range frames {

			// stop immediately if time exceeded
			if time.Since(start) >= duration {
				break
			}

			render(f)
			uiDelay(int(frameDelay / time.Millisecond))
		}
	}

	clearScreen()

	title := center(fmt.Sprintf("%s%sAhMyth%s", bold, blue, reset), width)
	sub := center("Android Remote Administration Tool", width)

	fmt.Println("\n\n\n" + title)
	fmt.Println(sub)

	uiDelay(3000)
	clearScreen()
}

/*func ShowSplash() {
	width := getWidth()

	fmt.Print(hide)
	defer fmt.Print(show)

	frames := []string{"│", "┃", "█", "┃", "│"}

	for i := 0; i < 3; i++ {
		for _, f := range frames {
			clearScreen()

			line := center(fmt.Sprintf("%s%s%s", cyan, f, reset), width)
			fmt.Println("\n\n\n" + line)

			time.Sleep(80 * time.Millisecond)
		}
	}

	clearScreen()

	title := center(fmt.Sprintf("%s%sAhMyth%s", bold, blue, reset), width)
	sub := center("Android Remote Administration Tool", width)

	fmt.Println("\n\n\n" + title)
	fmt.Println(sub)

	time.Sleep(1500 * time.Millisecond)
	clearScreen()
}
*/
