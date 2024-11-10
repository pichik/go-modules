package request

import (
	"fmt"
	"os"
	"time"

	"github.com/pichik/go-modules/misc"
)

func PrintUrl(requestData misc.RequestData, save bool) {
	var redirect string
	color := misc.White

	switch requestData.ResponseStatus {
	case 0:
		color = misc.Gray
	case 200, 204:
		color = misc.Green
	case 301, 303, 302:
		redirect = fmt.Sprintf(" ->%s %s", misc.White, requestData.ResponseHeaders.Get("Location"))
		color = misc.Yellow
	case 400, 500, 501, 422:
		color = misc.Blue
	case 403, 401:
		color = misc.Purple
	case 404, 429, 410:
		color = misc.Red
	case 405:
		color = misc.Orange
	default:
		color = misc.White
	}
	method := "[" + requestData.Method + "]"
	contentType := fmt.Sprintf("[%s(%d)]", requestData.ResponseContentType, requestData.ResponseContentLength)

	formattedOutput := fmt.Sprintf("\r%s%9s [%.3d] %23s [ %s ] %s%s", color, method, requestData.ResponseStatus, contentType, misc.BuildUrl(requestData.ParsedUrl, "1234"), redirect, misc.White)

	fmt.Printf("%s\n", formattedOutput)
	// ResultOutput(requestData, formattedOutput)
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
	fmt.Fprintf(os.Stderr, "%sProgress: [(%d/%d)/%d | %s/%s]", "\r", Resulted, curr, max, misc.ReadableTime(seconds), misc.ReadableTime(estimatedTime))
}

func SetSartTime() {
	startTime = time.Now()
}
