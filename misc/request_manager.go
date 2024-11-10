package request

import (
	"io"
	"net/http"
)

type RequestData struct {
	ParsedUrl             ParsedUrl
	RequestBody           string
	ReqBody               io.Reader
	Method                string
	Headers               map[string]string
	Cookies               []http.Cookie
	ResponseBody          string
	ResponseBodyBytes     []byte
	ResponseStatus        int
	ResponseContentType   string
	ResponseContentLength int
	ResponseHeaders       http.Header
	Error                 error
	Retries               int
}
