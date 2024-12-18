package data

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/pichik/go-modules/misc"
)

var backslashreg = regexp.MustCompile(`\\([^a-zA-Z0-9]|$)`)

func ReadResponse(requestData *misc.RequestData, res *http.Response) error {

	var reader io.Reader
	var err error

	// Check if response is compressed with gzip
	if res.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			requestData.Error = err
			misc.PrintError("Gzip read", err)
			return err
		}
	} else {
		reader = res.Body
	}

	var buffer bytes.Buffer

	io.Copy(&buffer, reader)

	if err != nil {
		return err
	}
	requestData.ResponseContentLength = len(buffer.Bytes())
	requestData.ResponseBody = string(buffer.Bytes())
	requestData.ResponseBodyBytes = buffer.Bytes()
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x3c", "<", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x3e", ">", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x3d", "=", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x22", "\"", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x27", "'", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x26", "&", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "&amp;", "&", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\t", " ", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\n", " ", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\r", " ", -1)
	// requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\/", "/", -1)
	// Remove all backslashes that are not followed by alphanum characters
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\\\", "\\", -1)
	requestData.ResponseBody = backslashreg.ReplaceAllString(requestData.ResponseBody, "")

	requestData.ResponseStatus = res.StatusCode

	contentreg := regexp.MustCompile(`/\s*(.*?)\s*(;|$)`)

	if contentType := res.Header.Get("content-type"); contentType != "" {
		filteredContentType := contentreg.FindStringSubmatch(contentType)
		requestData.ResponseContentType = filteredContentType[1]
	}
	requestData.ResponseHeaders = res.Header

	return nil
}
