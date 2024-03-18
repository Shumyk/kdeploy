package cmd

import (
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
		RedStderr("Error:", err)
	}
}

func Goodbye(s ...any) {
	PurpleStout(s)
	os.Exit(0)
}

func ErrorCheck(err error, msg ...any) {
	if err != nil {
		Error(msg, err)
	}
}

func Error(s ...any) {
	RedStderr(s)
	os.Exit(1)
}
