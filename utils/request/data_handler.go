package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pichik/go-modules/output"
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
	output.Write(layout, incompleteUrlsFile, UrlDir)
	output.Write(layout, completeUrlsFile, UrlDir)
	for _, d := range output.DataScrapRegex() {
		if d.Name != "urls" {
			output.Write(layout, d.Name, DataDir)
		}
	}
}

func DataOutput(foundData []output.ScrapData, completeUrls []string, incompleteUrls []string) {
	var new []string
	if toolDir == "crawler/" {

		// var urls []string

		for _, d := range foundData {

			if d.Name != "urls" {
				new = output.Anew(d.Results, allFile, DataDir, true)
				output.WriteAll(new, d.Name, DataDir)
			}
			//  else {
			// 	urls = append(urls, d.Results...)

			// }
		}

		// Unique(&urls)
		// output.Anew(urls, allFile, UrlDir, true)

		output.Anew(incompleteUrls, incompleteUrlsFile, UrlDir, true)
		output.Anew(completeUrls, completeUrlsFile, UrlDir, true)
	} else {

		for _, d := range foundData {
			if d.Name != "urls" {
				new = output.Anew(d.Results, allFile, DataDir, false)
				output.WriteAll(new, d.Name, DataDir+toolDir)
			} else {
				new = output.Anew(d.Results, allFile, UrlDir, false)
				output.WriteAll(new, allFile, UrlDir+toolDir)
			}
		}
		new = output.Anew(incompleteUrls, incompleteUrlsFile, UrlDir, false)
		output.WriteAll(new, incompleteUrlsFile, UrlDir+toolDir)
		new = output.Anew(completeUrls, completeUrlsFile, UrlDir, false)
		output.WriteAll(new, completeUrlsFile, UrlDir+toolDir)
	}

}

func CorsOutput(responseHeaders http.Header, currentUrl string) {
	allowOrigin := responseHeaders.Get("Access-Control-Allow-Origin")
	allowCreds := responseHeaders.Get("Access-Control-Allow-Credentials")

	if allowOrigin != "" && allowCreds != "" {
		text := fmt.Sprintf("\033[32m%s\n\t\033[33mAccess-Control-Allow-Origin: \t\t%s\n\tAccess-Control-Allow-Credentials:\t%s\n\033[0m", currentUrl, allowOrigin, allowCreds)
		output.Write(text, "cors", DataDir)
	}
}

func ResponseOutput(requestData RequestData) {

	dir := ResponseDir + toolDir

	queries, _ := json.MarshalIndent(requestData.ParsedUrl.Queries, "", "\t")
	headers, _ := json.MarshalIndent(requestData.Headers, "", "\t")

	output.Write(fmt.Sprintf("\033[33m%7s-----------------------------------------------------------------------------------------", ""), BuildUrl(requestData.ParsedUrl, "23"), dir)
	output.Write(fmt.Sprintf("%7s: %s\n%7s: %d\n%7s: %s\n%7s: %s\n%7s: %s", "Method", requestData.Method, "Status", requestData.ResponseStatus, "Query", string(queries), "Headers", string(headers), "Body", requestData.RequestBody), BuildUrl(requestData.ParsedUrl, "23"), dir)
	output.Write(fmt.Sprintf("%7s-----------------------------------------------------------------------------------------\033[0m", ""), BuildUrl(requestData.ParsedUrl, "23"), dir)
	output.Write(requestData.ResponseBody, BuildUrl(requestData.ParsedUrl, "23"), dir)
}

func ResultOutput(requestData RequestData, formattedOutput string) {
	if output.ExtensionPass(requestData.ParsedUrl.Extension) {
		output.Write(formattedOutput, resultsFile, toolDir)
	}
}

func AddToTested(url string) {
	output.Write(url, testedFile, toolDir)
	output.RemoveLine(url, queueFile, toolDir)
}

func FillQueue(urls []string, ignoreTested bool) {
	Unique(&urls)
	urlsToTest := []string{}

	if ignoreTested {
		urlsToTest = urls
		output.RemoveFile(toolDir + queueFile)
	} else {
		urlsToTest = output.Anew(urls, testedFile, toolDir, false)
	}
	output.Anew(urlsToTest, queueFile, toolDir, true)
}

func CustomOutput(text string, filename string) {
	output.Write(text, filename, toolDir)
}
func CustomOutputs(text []string, filename string) {
	output.Anew(text, filename, toolDir, true)
}

func ReadQueue() []string {
	content, _ := output.Read(toolDir + queueFile)
	return content
}
