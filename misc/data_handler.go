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
const incompleteUrlsFile = "incomplete-urls"

var toolDir = string("unknown/")

func SetToolDir(dir string) {
	toolDir = dir + "/"
}

func DataSeparator() {
	//Separate previous found data from current
	layout := fmt.Sprintf("-----------------------------------------%s", time.Now().Format("[02.01.2006][15:04:05]"))
	Write(layout, incompleteUrlsFile, UrlDir)
	Write(layout, completeUrlsFile, UrlDir)
	for _, d := range DataScrapRegex() {
		if d.Name != "urls" {
			Write(layout, d.Name, DataDir)
		}
	}
}

func DataOutput(foundData []ScrapData, completeUrls []string, incompleteUrls []string) {
	var new []string
	if toolDir == "crawler/" {

		// var urls []string

		for _, d := range foundData {

			if d.Name != "urls" {
				new = Anew(d.Results, allFile, DataDir, true)
				WriteAll(new, d.Name, DataDir)
			}
			//  else {
			// 	urls = append(urls, d.Results...)

			// }
		}

		// Unique(&urls)
		// Anew(urls, allFile, UrlDir, true)

		Anew(incompleteUrls, incompleteUrlsFile, UrlDir, true)
		Anew(completeUrls, completeUrlsFile, UrlDir, true)
	} else {

		for _, d := range foundData {
			if d.Name != "urls" {
				new = Anew(d.Results, allFile, DataDir, false)
				WriteAll(new, d.Name, DataDir+toolDir)
			} else {
				new = Anew(d.Results, allFile, UrlDir, false)
				WriteAll(new, allFile, UrlDir+toolDir)
			}
		}
		new = Anew(incompleteUrls, incompleteUrlsFile, UrlDir, false)
		WriteAll(new, incompleteUrlsFile, UrlDir+toolDir)
		new = Anew(completeUrls, completeUrlsFile, UrlDir, false)
		WriteAll(new, completeUrlsFile, UrlDir+toolDir)
	}

}

func CorsOutput(responseHeaders http.Header, currentUrl string) {
	allowOrigin := responseHeaders.Get("Access-Control-Allow-Origin")
	allowCreds := responseHeaders.Get("Access-Control-Allow-Credentials")

	if allowOrigin != "" && allowCreds != "" {
		text := fmt.Sprintf("\033[32m%s\n\t\033[33mAccess-Control-Allow-Origin: \t\t%s\n\tAccess-Control-Allow-Credentials:\t%s\n\033[0m", currentUrl, allowOrigin, allowCreds)
		Write(text, "cors", DataDir)
	}
}

func ResponseOutput(requestData RequestData) {

	dir := ResponseDir + toolDir

	queries, _ := json.MarshalIndent(requestData.ParsedUrl.Queries, "", "\t")
	headers, _ := json.MarshalIndent(requestData.Headers, "", "\t")

	Write(fmt.Sprintf("\033[33m%7s-----------------------------------------------------------------------------------------", ""), BuildUrl(requestData.ParsedUrl, "23"), dir)
	Write(fmt.Sprintf("%7s: %s\n%7s: %d\n%7s: %s\n%7s: %s\n%7s: %s", "Method", requestData.Method, "Status", requestData.ResponseStatus, "Query", string(queries), "Headers", string(headers), "Body", requestData.RequestBody), BuildUrl(requestData.ParsedUrl, "23"), dir)
	Write(fmt.Sprintf("%7s-----------------------------------------------------------------------------------------\033[0m", ""), BuildUrl(requestData.ParsedUrl, "23"), dir)
	Write(requestData.ResponseBody, BuildUrl(requestData.ParsedUrl, "23"), dir)
}

func ResultOutput(extension string, formattedOutput string) {
	if ExtensionPass(extension) {
		Write(formattedOutput, resultsFile, toolDir)
	}
}

func AddToTested(url string) {
	Write(url, testedFile, toolDir)
	RemoveLine(url, queueFile, toolDir)
}

func FillQueue(urls []string, ignoreTested bool) {
	Unique(&urls)
	urlsToTest := []string{}

	if ignoreTested {
		urlsToTest = urls
		RemoveFile(toolDir + queueFile)
	} else {
		urlsToTest = Anew(urls, testedFile, toolDir, false)
	}
	Anew(urlsToTest, queueFile, toolDir, true)
}

func CustomOutput(text string, filename string) {
	Write(text, filename, toolDir)
}
func CustomOutputs(text []string, filename string) {
	Anew(text, filename, toolDir, true)
}

func ReadQueue() []string {
	content, _ := Read(toolDir + queueFile)
	return content
}
