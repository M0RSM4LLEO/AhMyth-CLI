package cli  
  
import (  
	"fmt"  
	"time"  
)  
  
func ShowSplash() {  
	HideCursor()  
	defer ShowCursor()  
  
	frames := []string{"│", "┃", "█", "┃", "│"}  
  
	render := func(f string) {  
		ClearScreen()  
		line := Center(fmt.Sprintf("%s%s%s", Cyan, f, Reset))  
		fmt.Println("\n\n\n" + line)  
	}  
  
	// Animation  
	duration := 3000 * time.Millisecond  
	frameDelay := 100 * time.Millisecond  
	start := time.Now()  
  
	for time.Since(start) < duration {  
		for _, f := range frames {  
			if time.Since(start) >= duration {  
				break  
			}  
			render(f)  
			Delay(int(frameDelay / time.Millisecond))  
		}  
	}  
  
	// Final title screen  
	ClearScreen()  
	title := Center(fmt.Sprintf("%s%sAhMyth%s", Bold, Blue, Reset))  
	subtitle := Center("Android Remote Administration Tool")  
  
	fmt.Println("\n\n\n" + title)  
	fmt.Println(subtitle)  
  
	Delay(3000)  
	ClearScreen()  
}
