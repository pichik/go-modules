package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/pichik/go-modules/misc"
	"github.com/pichik/go-modules/tool"
)

var cookieFlag, agentFlag, originFlag, bodyFlag, httpMethodFlag, proxyFlag string
var normalizeFlag, ForceHTTP2Flag bool
var timeoutFlag, retriesFlag int
var headersFlag tool.ArrayStringFlag
var client *http.Client
var timeout time.Duration

var RequestBase misc.RequestData
var requestCancel context.CancelFunc

func (util Request) SetupFlags() []tool.UtilData {
	var flags []tool.FlagData
	flags = append(flags,
		tool.FlagData{
			Name:        "C",
			Description: "Request Cookies",
			Required:    false,
			Def:         "",
			VarStr:      &cookieFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "B",
			Description: "Request body",
			Required:    false,
			Def:         "",
			VarStr:      &bodyFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "A",
			Description: "User agent",
			Required:    false,
			Def:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
			VarStr:      &agentFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "O",
			Description: "Origin header",
			Required:    false,
			Def:         "",
			VarStr:      &originFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "T",
			Description: "Response timeout",
			Required:    false,
			Def:         30,
			VarInt:      &timeoutFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "R",
			Description: "Number of request retries on errors",
			Required:    false,
			Def:         3,
			VarInt:      &retriesFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "N",
			Description: "Normalize endpoints",
			Required:    false,
			Def:         true,
			VarBool:     &normalizeFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "H",
			Description: "Request Headers, use one header per flag",
			Required:    false,
			Def:         tool.ArrayStringFlag{},
			VarAStr:     &headersFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "X",
			Description: "Http method",
			Required:    false,
			Def:         "GET",
			VarStr:      &httpMethodFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "P",
			Description: "Proxy URL (e.g., http://proxy.example.com:8080)",
			Required:    false,
			Def:         "",
			VarStr:      &proxyFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "http",
			Description: "Force http2",
			Required:    false,
			Def:         true,
			VarBool:     &ForceHTTP2Flag,
		})

	util.UtilData.Name = "HTTP Request"
	util.UtilData.FlagDatas = flags

	return []tool.UtilData{*util.UtilData}
}

func (util Request) SetupData() {

	misc.Normalize = normalizeFlag

	RequestBase = misc.RequestData{
		Method:  httpMethodFlag,
		Cookies: setCookies(cookieFlag),
		Headers: setHeaders(headersFlag),
	}
	timeout = time.Duration(timeoutFlag)

	RequestBase.RequestBody = bodyFlag
	RequestBase.ReqBody = bytes.NewBuffer([]byte(bodyFlag))

	setClient()
	interruptMonitor()
}

func CreateRequest(requestData *misc.RequestData) {
	fmt.Println("PICA MAAAAAAAAAAAAAT")
	url := requestData.ParsedUrl.Url

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	requestCancel = cancel

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, requestData.Method, url, requestData.ReqBody)
	if err != nil {
		if requestData.Retries < 3 {
			time.Sleep(1 * time.Second)
			requestData.Retries++
			CreateRequest(requestData)
			return
		}
		requestData.Error = err
		fmt.Println(err)
		return
	}
	req.Close = true

	for k, v := range requestData.Headers {
		req.Header.Add(k, v)
	}
	for _, v := range requestData.Cookies {
		req.AddCookie(&v)
	}
	client.CloseIdleConnections()

	res, err := client.Do(req)

	if err != nil {
		requestData.Error = err
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	// res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	requestData.ResponseContentLength = len(body)
	requestData.ResponseBody = string(body)
	requestData.ResponseBodyBytes = body
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x3c", "<", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x3e", ">", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x3d", "=", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x22", "\"", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x27", "'", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\x26", "&", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "&amp;", "&", -1)
	requestData.ResponseBody = strings.Replace(requestData.ResponseBody, "\\/", "/", -1)
	requestData.ResponseStatus = res.StatusCode

	contentreg := regexp.MustCompile(`/\s*(.*?)\s*(;|$)`)

	if contentType := res.Header.Get("content-type"); contentType != "" {
		filteredContentType := contentreg.FindStringSubmatch(contentType)
		requestData.ResponseContentType = filteredContentType[1]
	}
	requestData.ResponseHeaders = res.Header
}

func setHeaders(headers []string) map[string]string {
	headerRegex := regexp.MustCompile(`(^[^\:]*): ?(.*$)`)
	tempHeaders := map[string]string{}
	for _, v := range headers {
		cn := headerRegex.FindStringSubmatch(v)
		tempHeaders[cn[1]] = cn[2]

	}
	tempHeaders["User-Agent"] = agentFlag
	tempHeaders["Origin"] = originFlag
	return tempHeaders
}

func setCookies(cookie string) []http.Cookie {
	cookieName := regexp.MustCompile(`(^|;)\s*(.*?)\s*=`)
	cookieValue := regexp.MustCompile(`\=\s*(.*?)\s*(;|$)`)

	cn := cookieName.FindAllStringSubmatch(cookie, -1)
	cv := cookieValue.FindAllStringSubmatch(cookie, -1)
	cookies := []http.Cookie{}
	for _, f := range cn {
		cookies = append(cookies, http.Cookie{Name: f[2]})
	}
	for i, f := range cv {
		cookies[i].Value = f[1]
	}
	return cookies
}

func setClient() {
	var proxy func(*http.Request) (*url.URL, error)

	if proxyFlag != "" {
		parsedProxy, err := url.Parse(proxyFlag)
		if err != nil {
			fmt.Println("Invalid proxy URL:", err)
			return
		}
		// Use http.ProxyURL to create a proxy function
		proxy = http.ProxyURL(parsedProxy)
	}

	client = &http.Client{
		Timeout: time.Second * timeout,

		//Dont follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			ForceAttemptHTTP2: ForceHTTP2Flag, // Forces HTTP/2
			Dial: (&net.Dialer{
				Timeout: timeout * time.Second,
			}).Dial,
			DialContext: (&net.Dialer{
				Timeout:   timeout * time.Second,
				KeepAlive: 1 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: timeout * time.Second,
			DisableKeepAlives:   true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Renegotiation:      tls.RenegotiateOnceAsClient,
			},
			Proxy: proxy,
		},
	}
}

func interruptMonitor() {
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range sigChan {
			if requestCancel != nil {
				requestCancel()
			}
			//Add empty line
			fmt.Println()
			os.Exit(0)
		}
	}()
}
