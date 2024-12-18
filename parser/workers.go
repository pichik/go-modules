package parser

import (
	"github.com/pichik/go-modules/misc"
)

// Worker function for goroutines
func jsobSelector(k string, v *ParserData, text string, currentUrl *misc.ParsedUrl, parsedUrls *[]misc.ParsedUrl) {

	switch k {
	case "urls":
		if includeUrls {
			*parsedUrls = append(*parsedUrls, findUrls(text, v, currentUrl)...)
		}
	default:
		if includeGithubData {
			getData(text, v)
		}
	}

	deduplicate(&v.Results)

}
