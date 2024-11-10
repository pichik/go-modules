package misc

import (
	"fmt"
	"regexp"
)

var extensions *regexp.Regexp

var eof *regexp.Regexp
var timeout *regexp.Regexp

var scrapData []ScrapData

type ScrapData struct {
	Name        string
	RegexString string
	Regex       *regexp.Regexp
	RegexPart   int
	Results     []string
	Highlight   string
}

func EOF() *regexp.Regexp {
	return eof
}
func Timeout() *regexp.Regexp {
	return timeout
}

func DataScrapRegex() []ScrapData {
	return scrapData
}

func CompileFilters() {
	extensions = regexp.MustCompile(`(jpe?g|png|svg|css|gif|ico|woff2?|ttf)`)

	eof = regexp.MustCompile(`EOF$`)
	timeout = regexp.MustCompile(`.*dial tcp.*i/o timeout$`)

	LoadFilters()

	for _, filter := range FilterData {
		var sd ScrapData
		sd.Name = filter.Name
		sd.Highlight = filter.Highlight
		sd.RegexPart = filter.RegexPart
		sd.Regex = regexp.MustCompile(fmt.Sprintf(`%s`, filter.RegexString))
		sd.RegexString = filter.RegexString
		scrapData = append(scrapData, sd)
	}

}

// Is not uselless file (jpe?g|png|svg|css|gif|ico|woff2?|ttf)
func ExtensionPass(extension string) bool {
	if extensions.MatchString(extension) {
		return false
	}
	return true
}
