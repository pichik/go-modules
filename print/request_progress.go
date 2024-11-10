package print

import (
	"fmt"
	"os"
	"strings"
	"time"
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

func PrintUrl(requestData RequestData, save bool) {
	var redirect string
	color := White

	switch requestData.ResponseStatus {
	case 0:
		color = Gray
	case 200, 204:
		color = Green
	case 301, 303, 302:
		redirect = fmt.Sprintf(" ->%s %s", White, requestData.ResponseHeaders.Get("Location"))
		color = Yellow
	case 400, 500, 501, 422:
		color = Blue
	case 403, 401:
		color = Purple
	case 404, 429, 410:
		color = Red
	case 405:
		color = Orange
	default:
		color = White
	}
	method := "[" + requestData.Method + "]"
	contentType := fmt.Sprintf("[%s(%d)]", requestData.ResponseContentType, requestData.ResponseContentLength)

	formattedOutput := fmt.Sprintf("\r%s%9s [%.3d] %23s [ %s ] %s%s", color, method, requestData.ResponseStatus, contentType, BuildUrl(requestData.ParsedUrl, "1234"), redirect, White)

	fmt.Printf("%s\n", formattedOutput)
	ResultOutput(requestData, formattedOutput)
}

func Highlight(text string, hightlight string) string {
	if hightlight != "" {
		return strings.ReplaceAll(text, hightlight, fmt.Sprintf("%s%s%s", Green, hightlight, White))
	}
	return text
}

var startTime time.Time
var Resulted int

func PrintProgress(curr int, max int) {

	currentTime := time.Now()
	seconds := int(currentTime.Sub(startTime).Seconds())
	var persec int
	var estimatedTime int
	if seconds > 0 {
		persec = curr / seconds
		if persec > 0 {
			estimatedTime = max / persec
		}
	}
	fmt.Fprintf(os.Stderr, "%sProgress: [(%d/%d)/%d | %s/%s]", "\r", Resulted, curr, max, readableTime(seconds), readableTime(estimatedTime))
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

func SetSartTime() {
	startTime = time.Now()
}

func PrintDebug(text any) {
	fmt.Fprintf(os.Stderr, "%s\n", text)

}
