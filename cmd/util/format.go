package cmd

import (
	"fmt"
	"strings"
	"time"
)

const (
	// ImageOptionFormat      = 2006-01-02 15:04:05     7d639e...     tags
	ImageOptionFormat         = "%v" + Divider + "%v" + Divider + "%v"
	Divider            string = "     "
	Delimiter          string = ","
	DigestPrefix       string = "sha256:"
	FriendlyDateFormat string = "2006-01-02 15:04:05"
)

func FormatImageOption(date time.Time, digest string, tags ...string) string {
	return fmt.Sprintf(ImageOptionFormat, Date(date), TrimDigestPrefix(digest), JoinComma(tags))
}

func Date(t time.Time) string {
	return t.Format(FriendlyDateFormat)
}

func TrimDigestPrefix(digest string) string {
	return strings.TrimPrefix(digest, DigestPrefix)
}

func JoinComma(parts []string) string {
	return strings.Join(parts, Delimiter)
}

func AppendSemicolon(tag string) string {
	if len(tag) > 0 {
		return ":" + tag
	}
	return ""
}
