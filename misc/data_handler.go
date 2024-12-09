package misc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const UrlDir = "urls/"
const DataDir = "data/"
const ResponseDir = "responses/"

const testedFile = "tested"
const queueFile = "queue"
const resultsFile = "results"

const allFile = ".all"
const completeUrlsFile = "complete-urls"

var toolDir = string("unknown/")

func SetToolDir(dir string) {
	toolDir = dir + "/"
}

func DataSeparator() {
	//Separate previous found data from current
	layout := fmt.Sprintf("-----------------------------------------%s", time.Now().Format("[02.01.2006][15:04:05]"))
	Append(completeUrlsFile, UrlDir, layout)
	for _, d := range DataScrapRegex() {
		if d.Name != "urls" {
			Append(d.Name, DataDir, layout)
		}
	}
}

func DataOutput(name string, foundData []string, completeUrls []string) {
	var new []string

	Anew(completeUrlsFile, UrlDir, true, completeUrls...)

	new = Anew(allFile, DataDir, true, foundData...)
	Append(name, DataDir, new...)

}

func CorsOutput(responseHeaders http.Header, currentUrl string) {
	allowOrigin := responseHeaders.Get("Access-Control-Allow-Origin")
	allowCreds := responseHeaders.Get("Access-Control-Allow-Credentials")

	if allowOrigin != "" && allowCreds != "" {
		text := fmt.Sprintf("\033[32m%s\n\t\033[33mAccess-Control-Allow-Origin: \t\t%s\n\tAccess-Control-Allow-Credentials:\t%s\n\033[0m", currentUrl, allowOrigin, allowCreds)
		Append("cors", "", text)
	}
}

func ResponseOutput(requestData RequestData) {

	dirname := ResponseDir + toolDir
	filename := BuildUrl(requestData.ParsedUrl, "23")

	queries, _ := json.MarshalIndent(requestData.ParsedUrl.Queries, "", "\t")
	headers, _ := json.MarshalIndent(requestData.Headers, "", "\t")

	Append(filename, dirname, fmt.Sprintf("\033[33m%7s-----------------------------------------------------------------------------------------", ""))
	Append(filename, dirname, fmt.Sprintf("%7s: %s\n%7s: %d\n%7s: %s\n%7s: %s\n%7s: %s", "Method", requestData.Method, "Status", requestData.ResponseStatus, "Query", string(queries), "Headers", string(headers), "Body", requestData.RequestBody))
	Append(filename, dirname, fmt.Sprintf("%7s-----------------------------------------------------------------------------------------\033[0m", ""))
	Append(filename, dirname, requestData.ResponseBody)
}

func ResultOutput(formattedOutput string) {
	Append(resultsFile, toolDir, formattedOutput)
}

func AddToTested(url string) {
	Anew(testedFile, toolDir, true, url)
	RemoveLine(url, queueFile, toolDir)
}

func FillQueue(urls []string, ignoreTested bool) {
	Unique(&urls)
	urlsToTest := []string{}

	if ignoreTested {
		urlsToTest = urls
	} else {
		urlsToTest = Anew(testedFile, toolDir, false, urls...)
	}
	Anew(queueFile, toolDir, true, urlsToTest...)
}

func RemoveQueueFile() {
	RemoveFile(toolDir + queueFile)
}

func CustomOutput(text string, filename string) {
	Append(filename, toolDir, text)
}
func CustomOutputs(text []string, filename string) {
	Anew(filename, toolDir, true, text...)
}

func ReadQueue() []string {
	content, _ := Read(toolDir + queueFile)
	return content
}
