package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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
	cli.Delay(1000)

	ok, detail = checkGo()
	allOK = printCheck("Go", ok, detail) && allOK

	cli.PrintSeparator()
	cli.Delay(1000)

	ok, detail = checkApktool()
	allOK = printCheck("apktool", ok, detail) && allOK

	cli.PrintSeparator()
	cli.Delay(1000)

	if !allOK {
		fmt.Printf("%sOne or more prerequisites failed.%s\n", cli.Red, cli.Reset)
		fmt.Println("Please install missing dependencies before continuing.")
		cli.PrintSeparator()
		cli.Delay(1000)
		os.Exit(1)
	}

	fmt.Printf("%sSystem ready. Starting...%s\n", cli.Green, cli.Reset)
	cli.Delay(1200)

	cli.Run()
}

// Helper functions
func printCheck(label string, ok bool, detail string) bool {
	symbolText := "ERR"
	symbolColor := cli.Red

	if ok {
		symbolText = "OK"
		symbolColor = cli.Green
	}

	if cli.SupportsUnicodeGlyphs() {
		symbolText = "✖"
		symbolColor = cli.Red
		if ok {
			symbolText = "✔"
			symbolColor = cli.Green
		}
	}

	fmt.Printf("%s %-24s %s%s%s\n", symbolColor+symbolText+cli.Reset, label, cli.Green, detail, cli.Reset)
	return ok
}

func checkJava() (bool, string) {
	javaPath, err := exec.LookPath("java")
	if err != nil {
		return false, "not found in PATH"
	}

	output, err := runCommandOutput(javaPath, "-version")
	if err != nil {
		return false, "found at " + javaPath + ", but failed to run"
	}

	version, ok := parseJavaVersion(output)
	if !ok {
		return true, "detected at " + javaPath + " (version could not be parsed)"
	}

	if version < 11 {
		return false, "version " + strconv.Itoa(version) + " detected at " + javaPath + " (need 11+)"
	}

	return true, "version " + strconv.Itoa(version) + " detected at " + javaPath
}

func checkGo() (bool, string) {
	goPath, err := exec.LookPath("go")
	if err != nil {
		return false, "not found in PATH"
	}

	output, err := runCommandOutput(goPath, "version")
	if err != nil {
		return false, "found at " + goPath + ", but failed to run"
	}

	version, ok := parseGoVersion(output)
	if !ok {
		return true, "detected at " + goPath + " (version could not be parsed)"
	}

	return true, "version " + version + " detected at " + goPath
}

func checkApktool() (bool, string) {
	apktoolPath, err := exec.LookPath("apktool")
	if err != nil {
		return false, "not found in PATH"
	}

	for _, flag := range []string{"-version", "--version"} {
		output, _ := runShellCommand("apktool " + flag)
		output = strings.TrimSpace(output)

		if output == "" {
			continue
		}

		if version, ok := parseApktoolVersion(output); ok {
			return true, "version " + version + " detected at " + apktoolPath
		}

		// Fallback: if apktool printed something useful, show it rather than hiding it.
		return true, "version " + output + " detected at " + apktoolPath
	}

	return true, "detected at " + apktoolPath + " (version unavailable)"
}

func runShellCommand(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func runCommandOutput(command string, args ...string) (string, error) {
	cmd := buildCommand(command, args...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func buildCommand(command string, args ...string) *exec.Cmd {
	// On Windows, .bat/.cmd files need cmd /C.
	if runtime.GOOS == "windows" {
		ext := strings.ToLower(filepath.Ext(command))
		if ext == ".bat" || ext == ".cmd" {
			return exec.Command("cmd", append([]string{"/C", command}, args...)...)
		}
	}

	// On Linux / Unix / Termux / WSL etc, run the command directly.
	return exec.Command(command, args...)
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

func parseGoVersion(output string) (string, bool) {
	// Example: go version go1.22.4 windows/amd64
	re := regexp.MustCompile(`go version go([^\s]+)`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", false
	}
	return match[1], true
}

func parseApktoolVersion(output string) (string, bool) {
	re := regexp.MustCompile(`\b(\d+\.\d+(?:\.\d+)*)\b`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", false
	}
	return match[1], true
}
