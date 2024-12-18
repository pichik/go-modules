package request

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pichik/go-modules/misc"
	"github.com/pichik/go-modules/request/data"
	"github.com/pichik/go-modules/tool"
)

type Request struct {
	UtilData *tool.UtilData
}

var cookieFlag, agentFlag, originFlag, bodyFlag, httpMethodFlag, proxyFlag string
var NormalizeFlag, ForceHTTP2Flag bool
var timeoutFlag, retriesFlag int
var headersFlag tool.ArrayStringFlag

var client *http.Client

var RequestBase misc.RequestData

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
			Def:         60,
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
			VarBool:     &NormalizeFlag,
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
			Def:         false,
			VarBool:     &ForceHTTP2Flag,
		})

	util.UtilData.Name = "HTTP Request"
	util.UtilData.FlagDatas = flags

	return []tool.UtilData{*util.UtilData}
}

func (util Request) SetupData() {
	misc.Normalize = NormalizeFlag

	headersFlag = append(headersFlag, fmt.Sprintf("Origin: %s", originFlag))
	headersFlag = append(headersFlag, fmt.Sprintf("User-Agent: %s", agentFlag))

	RequestBase = misc.RequestData{
		Method:  httpMethodFlag,
		Cookies: data.SetCookies(cookieFlag),
		Headers: data.SetHeaders(headersFlag),
	}
	RequestBase.RequestBody = bodyFlag
	RequestBase.ReqBody = bytes.NewBuffer([]byte(bodyFlag))

	client = data.SetClient(proxyFlag, ForceHTTP2Flag, time.Duration(timeoutFlag))
	data.InterruptMonitor()
}

func CreateRequest(requestData *misc.RequestData) {
	client.CloseIdleConnections()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutFlag))
	data.RequestCancel = cancel
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, requestData.Method, requestData.ParsedUrl.Url, requestData.ReqBody)
	if err != nil {
		repeatRequest(requestData, err, "Request Error: ")
		cancel()
		return
	}

	for k, v := range requestData.Headers {
		req.Header.Add(k, v)
	}
	for _, v := range requestData.Cookies {
		req.AddCookie(&v)
	}

	sendRequest(requestData, req)
	cancel()
}

func sendRequest(requestData *misc.RequestData, req *http.Request) {
	res, err := client.Do(req)
	if err != nil {
		repeatRequest(requestData, err, "Response Error: ")
		return
	}

	defer res.Body.Close()

	if err := data.ReadResponse(requestData, res); err != nil {
		repeatRequest(requestData, err, "Reading response Error: ")
	}
}

func repeatRequest(requestData *misc.RequestData, err error, errString string) {
	if requestData.Retries < 3 {
		time.Sleep(1 * time.Second)
		requestData.Retries++
		CreateRequest(requestData)
	}
	requestData.Error = err
	// misc.PrintError(errString, err)
}
