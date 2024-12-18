package data

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/pichik/go-modules/misc"
)

var RequestCancel context.CancelFunc

func SetCookies(cookie string) []http.Cookie {
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

func SetHeaders(headers []string) map[string]string {
	headerRegex := regexp.MustCompile(`(^[^\:]*): ?(.*$)`)
	tempHeaders := map[string]string{}
	for _, v := range headers {
		cn := headerRegex.FindStringSubmatch(v)
		tempHeaders[cn[1]] = cn[2]

	}
	return tempHeaders
}

func SetClient(proxyip string, forceHTTP2 bool, timeout time.Duration) *http.Client {
	var proxy func(*http.Request) (*url.URL, error)

	if proxyip != "" {
		parsedProxy, err := url.Parse(proxyip)
		if err != nil {
			misc.PrintError("Invalid proxy url", err)
			os.Exit(1)
		}
		// Use http.ProxyURL to create a proxy function
		proxy = http.ProxyURL(parsedProxy)
	}

	client := &http.Client{
		Timeout: time.Second * timeout,

		//Dont follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			ForceAttemptHTTP2: forceHTTP2, // Forces HTTP/2
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
	return client
}

func InterruptMonitor() {
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range sigChan {
			if RequestCancel != nil {
				RequestCancel()
			}
			//Add empty line
			fmt.Println()
			os.Exit(0)
		}
	}()
}
