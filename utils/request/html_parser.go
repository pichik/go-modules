package request

import (
	"fmt"
	"regexp"
	"sync"
)

func GetData(text string, parsedUrl *ParsedUrl) ([]output.ScrapData, []ParsedUrl, []ParsedUrl) {
	scrapedData := output.DataScrapRegex()
	retreivedDataChan := make(chan output.ScrapData, len(scrapedData))
	var wg sync.WaitGroup

	var completeUrls []ParsedUrl
	var incompleteUrls []ParsedUrl

	for _, s := range scrapedData {
		wg.Add(1)
		go func(s output.ScrapData, scraped chan output.ScrapData) {
			if s.Name == "urls" {
				findUrls(&s, text, parsedUrl, &completeUrls, &incompleteUrls)
			} else {
				for _, fd := range s.Regex.FindAllStringSubmatch(text, -1) {
					s.Results = append(s.Results, output.Highlight(fd[s.RegexPart], s.Highlight))
				}
				Unique(&s.Results)
			}
			scraped <- s
			wg.Done()
		}(s, retreivedDataChan)
	}

	wg.Wait()
	close(retreivedDataChan)

	for i := range scrapedData {
		scrapedData[i] = <-retreivedDataChan
	}

	UniqueUrls(&completeUrls)
	UniqueUrls(&incompleteUrls)

	return scrapedData, completeUrls, incompleteUrls
}

func findUrls(urlRegex *output.ScrapData, text string, currentUrl *ParsedUrl, completeUrls *[]ParsedUrl, incompleteUrls *[]ParsedUrl) {
	foundUrls := urlRegex.Regex.FindAllStringSubmatch(text, -1)

	for _, fu := range foundUrls {
		parsedUrl := ParseUrl(fu[urlRegex.RegexPart])

		if parsedUrl.Error != nil {
			urlRegex.Results = append(urlRegex.Results, parsedUrl.Url)
			if output.ExtensionPass(parsedUrl.Extension) {
				*incompleteUrls = append(*incompleteUrls, parsedUrl)
			}
			continue
		}

		//check if domain is missing in endpoints and add domain from current request
		if parsedUrl.Domain == "" && currentUrl != nil {
			parsedUrl.Domain = currentUrl.Domain
			RebuildUrl(&parsedUrl)
		}
		urlRegex.Results = append(urlRegex.Results, parsedUrl.Url)

		if output.ExtensionPass(parsedUrl.Extension) {
			*completeUrls = append(*completeUrls, parsedUrl)
		}

	}

	Unique(&urlRegex.Results)
}

func FindString(text string, toFind string, rangee int) []string {
	reg := fmt.Sprintf(`.{0,%d}%s.{0,%d}`, rangee, toFind, rangee)
	regex := regexp.MustCompile(reg)
	f := regex.FindAllString(text, -1)
	Unique(&f)
	return f
}

func Unique(slice *[]string) {
	seen := make(map[string]struct{})
	j := 0
	for _, s := range *slice {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		(*slice)[j] = s
		j++
	}
	*slice = (*slice)[:j]
}
