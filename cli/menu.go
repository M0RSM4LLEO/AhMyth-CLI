package cli  
  
import (  
	"bufio"  
	"fmt"  
	"os"  
	"strings"  
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
	ClearScreen()  
	HideCursor()  
	defer ShowCursor()  
  
	PrintLogo()  
	fmt.Println()  
  
	menuTitle := "Main Menu"  
	fmt.Printf("%s%s%s%s\n", Bold, Blue, menuTitle, Reset)  
	fmt.Printf("%s%s%s\n", Blue, strings.Repeat("—", len(menuTitle)), Reset)  
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
		fmt.Print("Enter choice: ")  
  
		input, _ := reader.ReadString('\n')  
		input = strings.TrimSpace(input)  
  
		switch input {  
		case "0":  
			ClearScreen()  
			return  
		case "1", "2":  
			fmt.Println("This feature is not implemented yet.")  
			waitForEnter(reader)  
		default:  
			fmt.Println("Invalid option!")  
			Delay(800)  
		}  
	}  
}
