package cmd

import (
	"fmt"
	"os"
)

func TerminateOnSigint(result string) {
	if len(result) == 0 {
		Goodbye("heh, ctrl+c combination was gently pressed. see you")
	}
}

func TerminateOnEmpty[T any](args []T, msg ...any) {
	if len(args) == 0 {
		Error(msg...)
	}
}

// Laugh just prints error message if present and ignores it
func Laugh(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func Goodbye(s ...any) {
	fmt.Println(purple(s))
	os.Exit(0)
}

func ErrorCheck(err error, msg ...any) {
	if err != nil {
		Error(err, msg)
	}
}

func Error(s ...any) {
	_, _ = fmt.Fprintln(os.Stderr, red(s...))
	os.Exit(1)
}
