package cmd

import (
	"os"
	"os/exec"
	"testing"
)

func TestTerminateOnSigint(t *testing.T) {
	if os.Getenv("IS_TEST") == "1" {
		TerminateOnSigint("")
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestTerminateOnSigint")
	cmd.Env = append(os.Environ(), "IS_TEST=1")
	err := cmd.Run()
	if err == nil {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 0", err)
}

func TestTerminateOnEmpty_EmptyInput(t *testing.T) {
	if os.Getenv("IS_TEST") == "1" {
		TerminateOnEmpty([]string{}, "error message")
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestTerminateOnEmpty")
	cmd.Env = append(os.Environ(), "IS_TEST=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && e.ExitCode() == 1 {
		return
	}
	t.Fatalf("process run with err %v, want exit status 1", err)
}

func TestTerminateOnEmpty_NonEmptyInput(t *testing.T) {
	TerminateOnEmpty([]string{"first", "second"}, "shouldn't appear")
}
