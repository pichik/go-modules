package parser

import (
	"regexp"
	"strings"
	"sync"

	"github.com/pichik/go-modules/misc"
)

var alphaNumRegex = regexp.MustCompile(`[A-Za-z0-9]{2,}`)
var containsNormalCharRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\.,!?'"` + "`" + `@#$%\^&*()_+\-=\[\]\{\};:\\<>\|/~]*$`)

func findUrls(text string, parserData *ParserData, currentUrl *misc.ParsedUrl) []misc.ParsedUrl {
	var completeUrls []misc.ParsedUrl
	var wg sync.WaitGroup // To wait for all goroutines to finish
	mu := sync.Mutex{}    // Mutex to protect shared resources (e.g., parserData.Results, completeUrls)

	for _, filter := range parserData.Filters {
		wg.Add(1) // Add a goroutine to the wait group

		go func(filter Filters) {
			defer wg.Done() // Mark this goroutine as done when it finishes

			// Run ag with the current filter's RegexString
			findings, err := runAg(text, filter.RegexString+".?")
			if err != nil {
				misc.PrintError("Regex failed", err)
				return
			}

			// Iterate over the findings in this goroutine
			for _, finding := range findings {

				// Check if result has some alphanumeric characters to filter out junk
				if !alphaNumRegex.MatchString(finding) || !strings.Contains(finding, "/") {
					continue
				}
				// Filter out another junk
				if strings.ContainsAny(finding, ";+") || !containsNormalCharRegex.MatchString(finding) {
					continue
				}
				// Match the finding against the filter's regex
				fixed := filter.Regex.FindStringSubmatch(finding)

				// Ensure valid match
				if filter.RegexPart >= len(fixed) || fixed[filter.RegexPart] == "" || strings.Contains(fixed[filter.RegexPart], " ") {
					continue
				}

				// Validate the URL and add it to completeUrls
				parsedUrl := misc.ParseUrl(fixed[filter.RegexPart])
				if !misc.ExtensionPass(parsedUrl.Protocol, parsedUrl.Extension) {
					continue
				}

				if parsedUrl.Error == nil {
					mu.Lock()
					completeUrls = append(completeUrls, parsedUrl)
					mu.Unlock()
				}
				// Lock the shared parserData.Results before modifying
				mu.Lock()
				parserData.Results = append(parserData.Results, fixed[filter.RegexPart])
				mu.Unlock()

			}
		}(filter) // Pass the filter as an argument to the goroutine
	}

	// Wait for all goroutines to finish before returning
	wg.Wait()

	misc.UniqueUrls(&completeUrls)

	return completeUrls
}

// Check if url pass url parser and filters
// func urlValidation(url string, currentUrl *misc.ParsedUrl) (misc.ParsedUrl, bool) {

// 	parsedUrl := misc.ParseUrl(url)

// 	if parsedUrl.Error != nil {
// 		return parsedUrl, false
// 	}

// 	//check if domain is missing in endpoints and add domain from current request
// 	if parsedUrl.Domain == "" && currentUrl != nil {
// 		parsedUrl.Domain = currentUrl.Domain
// 		misc.RebuildUrl(&parsedUrl)
// 	}

// 	return parsedUrl, true
// }
