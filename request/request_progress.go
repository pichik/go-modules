package request

import (
	"fmt"
	"os"
	"time"
)

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

func SetSartTime() {
	startTime = time.Now()
}
