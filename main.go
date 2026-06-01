package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ahmyth-cli/cli"
)

const (
	reset = "\033[0m"
	red   = "\033[31m"
	green = "\033[32m"
	blue  = "\033[34m"
	bold  = "\033[1m"
)

func stripAnsi(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}

func printHeader(text string) {
    visibleLen := len(stripAnsi(text))
    line := strings.Repeat("─", visibleLen)

    fmt.Println(text)
    fmt.Println(line)
}

func uiDelay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func clearTerminal() {
	fmt.Print("\033[2J\033[H")
}

func main() {
	// clear termanl and delay visuals 
	// for a nice clean look
	clearTerminal()
	fmt.Println()
	uiDelay(1000)

	//print the "Preflight Checks" string
	printHeader(fmt.Sprintf("%s%sPreflight Checks%s", bold, blue, reset))

	// add a 1 secone delay for neat visuals
	uiDelay(1000)

	allOK := true

	ok, detail := checkJava()
	allOK = printCheck("Java (OpenJDK 11+)", ok, detail) && allOK

	// needs a better seperator logic
	fmt.Println(strings.Repeat("─", 52))

	// add a 1 secone delay for neat vi>
        uiDelay(1000)

	ok, detail = checkApktool()
	allOK = printCheck("apktool", ok, detail) && allOK

	// needs better seperator logic
	fmt.Println(strings.Repeat("─", 52))

	// add a 1 secone delay for neat vi>
        uiDelay(1000)

	if !allOK {
		fmt.Println(fmt.Sprintf("%sOne or more prerequisites failed.%s", red, reset))

		// add a 1 secone delay for neat vi>
	        uiDelay(1000)

		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("%sSystem ready. Starting...%s", green, reset))
	time.Sleep(1200 * time.Millisecond)
	uiDelay(1000)

	cli.Run()
}

func printCheck(label string, ok bool, detail string) bool {
	symbol := fmt.Sprintf("%s✖%s", red, reset)
	status := red
	if ok {
		symbol = fmt.Sprintf("%s✔%s", green, reset)
		status = green
	}

	fmt.Printf("%s %-24s %s%s%s\n", symbol, label, status, detail, reset)
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
		return false, "found, but failed to run"
	}

	version, ok := parseJavaMajorVersion(string(output))
	if !ok {
		return false, "found, but version could not be determined"
	}

	if version < 11 {
		return false, fmt.Sprintf("version %d detected at %s (need 11+)", version, javaPath)
	}

	return true, fmt.Sprintf("version %d detected at %s", version, javaPath)
}

func checkApktool() (bool, string) {
	apktoolPath, err := exec.LookPath("apktool")
	if err != nil {
		return false, "not found in PATH"
	}

	cmd := exec.Command(apktoolPath, "-version")
	output, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(output))
	if err == nil && out != "" {
		return true, fmt.Sprintf("version %s detected at %s", out, apktoolPath)
	}

	shellCmd := exec.Command("sh", "-c", shellQuote(apktoolPath)+" -version")
	output, err = shellCmd.CombinedOutput()
	out = strings.TrimSpace(string(output))
	if err != nil && out == "" {
		return false, fmt.Sprintf("found at %s, but failed to execute", apktoolPath)
	}

	if out == "" {
		return true, fmt.Sprintf("installed at %s", apktoolPath)
	}

	return true, fmt.Sprintf("version %s detected at %s", out, apktoolPath)
}

func parseJavaMajorVersion(output string) (int, bool) {
	re := regexp.MustCompile(`version\s+"([^"]+)"`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return 0, false
	}

	raw := match[1]

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

	parts := strings.Split(raw, ".")
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, false
	}

	return n, true
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}
