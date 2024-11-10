package request

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

type ParsedUrl struct {
	Url            string
	Protocol       string
	Domain         string
	Path           string
	Extension      string
	NormalizedPath string
	ParsedPath     []string
	QueryString    string
	Fragments      string
	Queries        url.Values
	Error          error
}

var Normalize bool

func ParseUrl(u string) ParsedUrl {
	var parsedUrl ParsedUrl

	if strings.HasPrefix(u, "://") {
		u = "https" + u
	}

	pu, err := url.Parse(u)

	if err != nil {
		parsedUrl.Error = err
		parsedUrl.Url = u
		return parsedUrl
	}

	parsedUrl.Protocol = pu.Scheme
	if !strings.Contains(parsedUrl.Protocol, "http") {
		parsedUrl.Protocol = "https"
	}

	parsedUrl.Domain = pu.Hostname()
	//Use rawPath for maximum accuracy
	if pu.RawPath != "" {
		parsedUrl.Path = pu.RawPath
	} else {
		parsedUrl.Path = pu.EscapedPath()
	}

	if !strings.HasPrefix(parsedUrl.Path, "/") {
		parsedUrl.Path = "/" + parsedUrl.Path
	}

	if strings.Contains(parsedUrl.Path, "\\") && !strings.Contains(parsedUrl.Path, "\\u00") {
		parsedUrl.Path = strings.ReplaceAll(parsedUrl.Path, "\\", "")
	}

	parsedUrl.Extension = path.Ext(parsedUrl.Path)
	parsedUrl.QueryString = pu.RawQuery
	parsedUrl.Fragments = pu.Fragment
	parsedUrl.Queries = pu.Query()

	if len(parsedUrl.Path) > 1 {
		normalizeEndpoint(&parsedUrl)
	}

	RebuildUrl(&parsedUrl)

	if strings.ContainsAny(pu.RawPath, "[{(;%") || strings.Contains(parsedUrl.Path, "\\u00") {
		parsedUrl.Error = errors.New("Incomplete url")
	}

	return parsedUrl
}

func ParseUrls(urls []string) []ParsedUrl {
	parsedUrls := []ParsedUrl{}
	for _, u := range urls {
		parsedUrls = append(parsedUrls, ParseUrl(u))
	}
	return parsedUrls
}

// Strip url to a parts, also build complete parsed.url from parts (ex: if parsed.domain is update, build it so parsed.url is wiith correct domain)
// Use 1234 to determine which parts should url contains
// 1 - Protocol
// 2 - Domain
// 3 - Endpoint
// 4 - Query and Fragments
func BuildUrl(parsedUrl ParsedUrl, parts string) string {
	url := string("")

	for _, ch := range parts {
		if string(ch) == "1" && parsedUrl.Domain != "" {
			url = parsedUrl.Protocol + "://"
		}
		if string(ch) == "2" && parsedUrl.Domain != "" {
			url = url + parsedUrl.Domain
		}
		if string(ch) == "3" {
			if !Normalize && parsedUrl.Path != "" {
				url = url + parsedUrl.Path

			} else if Normalize && parsedUrl.NormalizedPath != "" {
				url = url + parsedUrl.NormalizedPath
			}
		}
		if string(ch) == "4" && parsedUrl.Extension != "js" {
			if parsedUrl.QueryString != "" {
				url = url + "?" + parsedUrl.QueryString
			}
			if parsedUrl.Fragments != "" {
				url = url + "#" + parsedUrl.Fragments
			}
		}
	}
	return url
}

func BuildUrls(urls []ParsedUrl, parts string) []string {
	strippedUrls := []string{}

	for _, parsedUrl := range urls {
		var url string
		for _, ch := range parts {
			if string(ch) == "1" && parsedUrl.Domain != "" {
				url = url + parsedUrl.Protocol + "://"
			}
			if string(ch) == "2" && parsedUrl.Domain != "" {
				url = url + parsedUrl.Domain
			}
			if string(ch) == "3" {
				if !Normalize && parsedUrl.Path != "" {
					url = url + parsedUrl.Path
				} else if Normalize && parsedUrl.NormalizedPath != "" {
					url = url + parsedUrl.NormalizedPath
				}
			}
			if string(ch) == "4" && parsedUrl.Extension != "js" {
				if parsedUrl.QueryString != "" {
					url = url + "?" + parsedUrl.QueryString
				}
				// if parsedUrl.Fragments != "" {
				// 	url = url + "#" + parsedUrl.Fragments
				// }
			}
		}
		strippedUrls = append(strippedUrls, url)
	}
	return strippedUrls
}

func RebuildUrl(parsedUrl *ParsedUrl) {
	parsedUrl.Url = ""

	if parsedUrl.Domain != "" {
		parsedUrl.Url = parsedUrl.Protocol + "://"
	}
	if parsedUrl.Domain != "" {
		parsedUrl.Url = parsedUrl.Url + parsedUrl.Domain
	}

	if !Normalize && parsedUrl.Path != "" {
		parsedUrl.Url = parsedUrl.Url + parsedUrl.Path

	} else if Normalize && parsedUrl.NormalizedPath != "" {
		parsedUrl.Url = parsedUrl.Url + parsedUrl.NormalizedPath

	}

	if parsedUrl.QueryString != "" {
		parsedUrl.Url = parsedUrl.Url + "?" + parsedUrl.QueryString
	}
	if parsedUrl.Fragments != "" {
		parsedUrl.Url = parsedUrl.Url + "#" + parsedUrl.Fragments
	}
}

func GetUrls(parsedUrls []ParsedUrl) []string {
	var values []string
	for _, url := range parsedUrls {
		values = append(values, url.Url)
	}
	return values
}

func UniqueUrls(parsedUrls *[]ParsedUrl) {
	seen := make(map[string]struct{})
	j := 0
	for _, parsedUrl := range *parsedUrls {
		if parsedUrl.Url == "" {
			continue
		}
		if _, ok := seen[parsedUrl.Url]; ok {
			continue
		}
		seen[parsedUrl.Url] = struct{}{}
		(*parsedUrls)[j] = parsedUrl
		j++
	}
	*parsedUrls = (*parsedUrls)[:j]
}

func normalizeEndpoint(parsedUrl *ParsedUrl) {
	parsedEndpoint := strings.Split(parsedUrl.Path, "/")

	var normalizedEndpoint string

	for i := len(parsedEndpoint) - 1; i > 0; i-- {
		if parsedEndpoint[i] == ".." {
			i--
			continue
		} else if parsedEndpoint[i] == "." || parsedEndpoint[i] == "" {
			continue
		}
		parsedUrl.ParsedPath = append([]string{parsedEndpoint[i]}, parsedUrl.ParsedPath...)
		normalizedEndpoint = "/" + parsedEndpoint[i] + normalizedEndpoint
	}
	parsedUrl.NormalizedPath = normalizedEndpoint
	RebuildUrl(parsedUrl)
}
