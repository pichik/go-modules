package request

import (
	"fmt"
	"os"
	"time"

	"github.com/pichik/go-modules/misc"
	"github.com/pichik/go-modules/output"
)

func PrintUrl(requestData misc.RequestData, save bool) {
	var redirect string
	color := output.White

	switch requestData.ResponseStatus {
	case 0:
		color = output.Gray
	case 200, 204:
		color = output.Green
	case 301, 303, 302:
		redirect = fmt.Sprintf(" ->%s %s", output.White, requestData.ResponseHeaders.Get("Location"))
		color = output.Yellow
	case 400, 500, 501, 422:
		color = output.Blue
	case 403, 401:
		color = output.Purple
	case 404, 429, 410:
		color = output.Red
	case 405:
		color = output.Orange
	default:
		color = output.White
	}
	method := "[" + requestData.Method + "]"
	contentType := fmt.Sprintf("[%s(%d)]", requestData.ResponseContentType, requestData.ResponseContentLength)

	formattedOutput := fmt.Sprintf("\r%s%9s [%.3d] %23s [ %s ] %s%s", color, method, requestData.ResponseStatus, contentType, misc.BuildUrl(requestData.ParsedUrl, "1234"), redirect, output.White)

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
	fmt.Fprintf(os.Stderr, "%sProgress: [(%d/%d)/%d | %s/%s]", "\r", Resulted, curr, max, output.ReadableTime(seconds), output.ReadableTime(estimatedTime))
}

func SetSartTime() {
	startTime = time.Now()
}
