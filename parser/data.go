package parser

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/pichik/go-modules/misc"
)

type Filters struct {
	Type        string
	Name        string
	RegexString string
	RegexPart   int
	Highlight   string
	Regex       *regexp.Regexp
}

var filterData []Filters

type ParserData struct {
	Filters []Filters
	Results []string
}

var parserDataTemplate map[string]ParserData

func loadFilters() {
	js, _ := misc.LoadGithubWordlist("filters.json")
	err := json.Unmarshal(js, &filterData)

	if err != nil {
		misc.PrintError("Unmarshaling filter from website", err)
	}

	compileFilters()
}

func compileFilters() {
	parserDataTemplate = make(map[string]ParserData)

	for _, filter := range filterData {
		filter.Regex = regexp.MustCompile(fmt.Sprintf(`%s`, filter.RegexString))

		// Check if the key exists in the map
		if _, exists := parserDataTemplate[filter.Type]; !exists {
			// Initialize a new ParserData entry
			parserDataTemplate[filter.Type] = ParserData{
				Filters: []Filters{filter},
				Results: []string{},
			}
		} else {
			// Append the filter to the existing entry
			parserData := parserDataTemplate[filter.Type]
			parserData.Filters = append(parserData.Filters, filter)
			parserDataTemplate[filter.Type] = parserData
		}
	}

}

func deduplicate(slice *[]string) {
	if len(*slice) < 2 {
		return
	}

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
