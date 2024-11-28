package functions

import (
	"fmt"
	"os"
	"strings"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
var Pink = "\033[38;2;255;105;180m" 

// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("%s%s%s\n", Red, fmt.Sprintf("error: %s", err), Reset)
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("%s%s%s\n", Cyan, fmt.Sprintf(format, args...), Reset)
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("%s%s%s\n", Yellow, fmt.Sprintf(format, args...), Reset)
}

// Success should be used to display success messages
func Success(message string) {
	fmt.Printf("%s%s%s\n", Green, message, Reset)
}

// Executing should be used to display execution messages
func Executing(message string) {
	fmt.Printf("%s%s%s\n", Gray, message, Reset)
}

// Command should be used to display command execution
func Command(message string) {
	fmt.Printf("%s%s%s\n", Pink, message, Reset)
}

// Attempting should be used to display attempts
func Attempting(message string) {
	fmt.Printf("%s%s%s...\n", Magenta, message, Reset)
}
