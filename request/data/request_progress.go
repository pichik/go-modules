package data

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pichik/go-modules/misc"
)

var startTime time.Time
var Resulted int
var CurrentThreads int

var ProgressCounter int
var ProgressMax int

func PrintUrl(requestData misc.RequestData, save bool) {
	var redirect string
	color := misc.White

	switch requestData.ResponseStatus {
	case 0:
		redirect = printError(requestData.Error.Error())
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

	if save {
		misc.ResultOutput(formattedOutput)
	}
}

func PrintProgress() {

	currentTime := time.Now()
	seconds := int(currentTime.Sub(startTime).Seconds())
	var persec int
	var estimatedTime int

	//Check if more requests are processed in one second, or one request is processed in more seconds and calculate maximum time
	if seconds > 0 {
		persec = ProgressCounter / seconds
		if persec > 0 {
			estimatedTime = ProgressMax / persec
		} else {
			persec = seconds / ProgressCounter
			estimatedTime = ProgressMax * persec
		}
	}
	fmt.Fprintf(os.Stderr, "\rProgress: [(%d/%d)/%d | %s/%s] [%d r/s]", Resulted, ProgressCounter, ProgressMax, misc.ReadableTime(seconds), misc.ReadableTime(estimatedTime), CurrentThreads)
}

func SetSartTime() {
	startTime = time.Now()
}

func printError(err string) string {
	// Split by colon and get the last part
	parts := strings.Split(err, ":")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[len(parts)-1])
	}
	return err
}
