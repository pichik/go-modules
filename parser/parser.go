package parser

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/pichik/go-modules/misc"
	"github.com/pichik/go-modules/tool"
)

type Parser struct {
	UtilData *tool.UtilData
}

var includeUrls, includeGithubData = true, true

func (util Parser) SetupFlags() []tool.UtilData {
	var flags []tool.FlagData

	util.UtilData.Name = "Scrape custom data and urls"
	util.UtilData.FlagDatas = flags

	var ut []tool.UtilData
	ut = append(ut, Filter{UtilData: &tool.UtilData{}}.SetupFlags()...)
	ut = append(ut, *util.UtilData)

	return ut
}

func (util Parser) SetupData() {
	loadFilters()

	var ut []tool.IUtil
	ut = append(ut, Filter{})
	for _, u := range ut {
		u.SetupData()
	}
}

// Specify which data to get only when parsing recursively directory, as it may take too long to get all of them
func ParseDirectory(currentUrl *misc.ParsedUrl, urls bool, githubData bool) (map[string]ParserData, []misc.ParsedUrl) {
	includeUrls = urls
	includeGithubData = githubData

	recursiveCurrentDir = true
	return startParsing("", currentUrl)
}

func ParseText(text string, currentUrl *misc.ParsedUrl) (map[string]ParserData, []misc.ParsedUrl) {
	recursiveCurrentDir = false
	return startParsing(text, currentUrl)
}

func ParseTextWithRange(txt string, rangee int) {
}

func startParsing(text string, currentUrl *misc.ParsedUrl) (map[string]ParserData, []misc.ParsedUrl) {
	var parsedUrls []misc.ParsedUrl
	var wg sync.WaitGroup
	var mu sync.Mutex

	parserData := make(map[string]ParserData)

	// Channel to limit the number of concurrent goroutines
	sem := make(chan struct{}, 20) // Limit to 5 goroutines

	for k, v := range parserDataTemplate {
		// Acquire a token to allow starting a new goroutine
		sem <- struct{}{}
		wg.Add(1)

		// Start the goroutine
		go func(k string, v *ParserData) {
			defer wg.Done()
			defer func() {
				// Release the token when the goroutine finishes
				<-sem
			}()
			jsobSelector(k, v, text, currentUrl, &parsedUrls)

			mu.Lock()
			parserData[k] = *v
			mu.Unlock()
		}(k, &v)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return parserData, parsedUrls
}

func FindCustomString(text string, toFind string, rangee int) []string {
	reg := fmt.Sprintf(`.{0,%d}%s.{0,%d}`, rangee, toFind, rangee)
	regex := regexp.MustCompile(reg)
	f := regex.FindAllString(text, -1)
	deduplicate(&f)
	return f
}

func MergeData(maps ...map[string]ParserData) map[string]ParserData {
	result := make(map[string]ParserData)
	// Iterate over all maps
	for _, m := range maps {
		for key, value := range m {
			// Merge data (overwrite if key exists in result)
			if existing, exists := result[key]; exists {
				// You can merge the results here as needed, for example:
				existing.Results = append(existing.Results, value.Results...)
				deduplicate(&existing.Results)
				result[key] = existing
			} else {
				result[key] = value
			}
		}
	}
	return result
}
