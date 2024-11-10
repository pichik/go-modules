package output

import (
	"fmt"
	"os"
	"strings"
)

type Color int64

const (
	White Color = iota
	Gray
	Yellow
	Green
	Red
	Orange
	Blue
	Purple
)

func (c Color) String() string {
	switch c {
	case White:
		return "\033[0m"
	case Green:
		return "\033[32m"
	case Red:
		return "\033[31m"
	case Orange:
		return "\033[38;5;208m"
	case Purple:
		return "\033[35m"
	case Yellow:
		return "\033[33m"
	case Gray:
		return "\033[90m"
	case Blue:
		return "\033[36m"
	}
	return "unknown"
}

func Highlight(text string, hightlight string) string {
	if hightlight != "" {
		return strings.ReplaceAll(text, hightlight, fmt.Sprintf("%s%s%s", Green, hightlight, White))
	}
	return text
}

func readableTime(seconds int) string {
	hours := int(seconds / (60 * 60) % 24)
	minutes := int(seconds/60) % 60

	var h string
	var m string

	if hours > 0 {
		h = fmt.Sprintf("%2dh", hours)
	}
	if minutes > 0 {
		m = fmt.Sprintf("%2dm", minutes)
	}
	seconds = int(seconds % 60)
	return fmt.Sprintf("%s%s%2ds", h, m, seconds)
}

func PrintDebug(text any) {
	fmt.Fprintf(os.Stderr, "%s\n", text)

}
