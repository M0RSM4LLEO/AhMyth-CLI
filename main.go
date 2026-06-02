package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"ahmyth-cli/cli"
)

func main() {
	cli.ClearScreen()
	fmt.Println()
	cli.Delay(1000)

	cli.PrintHeader(fmt.Sprintf("%s%sPreflight Checks%s", cli.Bold, cli.Blue, cli.Reset))  
	cli.Delay(1000)

	allOK := true

	ok, detail := checkJava()
	allOK = printCheck("Java (OpenJDK 11+)", ok, detail) && allOK  

	cli.PrintSeparator()
	//fmt.Println(strings.Repeat("─", 52))  
	cli.Delay(1000)

	ok, detail = checkApktool()
	allOK = printCheck("apktool", ok, detail) && allOK  

	cli.PrintSeparator()
	//fmt.Println(strings.Repeat("─", 52))  
        cli.Delay(1000)

	if !allOK {
		fmt.Printf("%sOne or more prerequisites failed.%s\n", cli.Red, cli.Reset)  
		fmt.Println("Please install missing dependencies before continuing.")  
		fmt.Println(strings.Repeat("─", 52))  
	        cli.Delay(1000)
		os.Exit(1)
	}

	fmt.Printf("%sSystem ready. Starting...%s\n", cli.Green, cli.Reset)  
	cli.Delay(1200)

	cli.Run()
}

// Helper functions
func printCheck(label string, ok bool, detail string) bool {  
	symbol := fmt.Sprintf("%s✖%s", cli.Red, cli.Reset)  
	status := cli.Red
	if ok {
		symbol = fmt.Sprintf("%s✔%s", cli.Green, cli.Reset)  
		status = cli.Green
	}

	fmt.Printf("%s %-24s %s%s%s\n", symbol, label, status, detail, cli.Reset)  
	return ok
}

func checkJava() (bool, string) {
	javaPath, err := exec.LookPath("java")
	if err != nil {
		return false, "not found in PATH"
	}

	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, "found at " + javaPath + ", but failed to run"
	}

	version, ok := parseJavaVersion(string(output))
	if !ok {
		return true, "detected at " + javaPath + " (version could not be parsed)"
	}

	if version < 11 {
		return false, "version " + strconv.Itoa(version) + " detected at " + javaPath + " (need 11+)"
	}

	return true, "version " + strconv.Itoa(version) + " detected at " + javaPath
}

func checkApktool() (bool, string) {
	apktoolPath, err := exec.LookPath("apktool")
	if err != nil {
		return false, "not found in PATH"
	}

	// Run through shell like the terminal does
	cmd := exec.Command("sh", "-c", "apktool -version")
	output, _ := cmd.CombinedOutput()

	version := strings.TrimSpace(string(output))

	// Fallback: some versions use --version
	if version == "" {
		cmd = exec.Command("sh", "-c", "apktool --version")
		output, _ = cmd.CombinedOutput()
		version = strings.TrimSpace(string(output))
	}

	// If we got any output at all, treat it as success
	if version != "" {
		return true, "version " + version + " detected at " + apktoolPath
	}

	// Last fallback: executable exists, but version failed
	return true, "detected at " + apktoolPath + " (version unavailable)"
}

func parseJavaVersion(output string) (int, bool) {
	re := regexp.MustCompile(`version\s+"([^"]+)"`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return 0, false
	}

	raw := match[1]

	// Old Java format: 1.8.0_xxx
	if strings.HasPrefix(raw, "1.") {
		parts := strings.Split(raw, ".")
		if len(parts) < 2 {
			return 0, false
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, false
		}
		return n, true
	}

	// New Java format: 11, 17.0.2, 21.0.1, etc.
	parts := strings.Split(raw, ".")
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, false
	}
	return n, true
}
