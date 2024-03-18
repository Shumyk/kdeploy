package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
)

const (
	LengthInfoLine      = 19
	HeaderCurrentImage  = "CURRENT IMAGE"
	HeaderDeployedImage = "DEPLOYED IMAGE"
)

var (
	isDebug = false

	termWidth, _, _ = term.GetSize(int(os.Stdin.Fd()))
	header          = color.New(color.Bold, color.BgHiGreen).SprintFunc()
	green           = color.New(color.Bold, color.FgHiGreen).SprintFunc()
	purple          = color.New(color.Bold, color.FgMagenta).SprintFunc()
	cyan            = color.New(color.Bold, color.FgCyan).SprintFunc()
	red             = color.New(color.Bold, color.FgHiRed).SprintFunc()
)

func SetDebugMode(isDebugMode bool) {
	isDebug = isDebugMode
}

func Debug(msg ...any) {
	if isDebug {
		CyanStout("[DEBUG]", msg)
	}
}

// Laugh just prints error message if present and ignores it
func Laugh(err error, msg ...any) {
	if err != nil {
		RedStderr("Error: ", err, msg)
	}
}

func DashLine() {
	fmt.Printf("%s", strings.Repeat("-", termWidth))
}

func PrintEnvironmentInfo(service, namespace string) {
	printInfoBlock(
		"ENVIRONMENT",
		EntryOf("service", service),
		EntryOf("namespace", namespace),
	)
}

func PrintImageInfo(header, tag, digest string) {
	printInfoBlock(
		header,
		EntryOf("tag", tag),
		EntryOf("digest", digest),
	)
}

func printInfoBlock(header string, lines ...Entry) {
	wrapHeader(buildHeaderLine(header))
	for _, line := range lines {
		fmt.Println(buildInfoLine(line.Key, green(line.Value)))
	}
	DashLine()
}

func wrapHeader(title string) {
	DashLine()
	line := withTrailingWhitespaces(title)
	fmt.Println(header(line))
	DashLine()
}

func withTrailingWhitespaces(prefix string) string {
	trailingWhitespaces := strings.Repeat(" ", termWidth-len(prefix))
	return fmt.Sprintf("%v%v", prefix, trailingWhitespaces)
}

// terminal indentation helpers
// ↓↓↓						↓↓↓

func buildLine(msg, suffix string) string {
	prefix := fmt.Sprintf("|   %v", msg)
	freeSpace := LengthInfoLine - len(prefix)
	spaces := strings.Repeat(" ", freeSpace)
	return fmt.Sprintf("%v%v%v", prefix, spaces, suffix)
}

func buildHeaderLine(header string) string {
	return buildLine(header, "|")
}

func buildInfoLine(key, value string) string {
	suffix := fmt.Sprintf(":  %v", value)
	return buildLine(key, suffix)
}

// ↑↑↑				      	    ↑↑↑
// end terminal indentation helpers

// colors
// ↓↓↓						↓↓↓

func PurpleStout(msg ...any) {
	_, _ = fmt.Println(purple(msg...))
}

func CyanStout(msg ...any) {
	_, _ = fmt.Println(cyan(msg...))
}

func BoringStderr(msg ...any) {
	_, _ = fmt.Fprintln(os.Stderr, msg...)
}

func RedStderr(msg ...any) {
	_, _ = fmt.Fprintln(os.Stderr, red(msg...))
}

// ↑↑↑				      	    ↑↑↑
// colors
