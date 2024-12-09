package misc

import (
	"regexp"
)

var extensions *regexp.Regexp

// var eof *regexp.Regexp
// var timeout *regexp.Regexp
// var refused *regexp.Regexp
var repeatRequestTriggers *regexp.Regexp

var scrapData []ScrapData

type ScrapData struct {
	Name        string
	RegexString string
	Regex       *regexp.Regexp
	RegexPart   int
	Results     []string
	Highlight   string
}

// func EOF() *regexp.Regexp {
// 	return eof
// }
// func Timeout() *regexp.Regexp {
// 	return timeout
// }
// func Refused() *regexp.Regexp {
// 	return refused
// }

func RepeatRequestTriggers() *regexp.Regexp {
	return repeatRequestTriggers
}

func DataScrapRegex() []ScrapData {
	return scrapData
}

func CompileFilters() {
	extensions = regexp.MustCompile(`(jpe?g|png|svg|css|gif|ico|woff2?|ttf)`)

	// eof = regexp.MustCompile(`EOF$`)
	// timeout = regexp.MustCompile(`.*dial tcp.*i/o timeout$`)
	// refused = regexp.MustCompile(`connect: connection refused$`)

	repeatRequestTriggers = regexp.MustCompile(`(context deadline exceeded|EOF|dial tcp.*i/o timeout|connection refused)$`)

	// LoadFilters()

	// for _, filter := range FilterData {
	// 	var sd ScrapData
	// 	sd.Name = filter.Name
	// 	sd.Highlight = filter.Highlight
	// 	sd.RegexPart = filter.RegexPart
	// 	sd.Regex = regexp.MustCompile(fmt.Sprintf(`%s`, filter.RegexString))
	// 	sd.RegexString = filter.RegexString
	// 	scrapData = append(scrapData, sd)
	// }

}

// Is not uselless file (jpe?g|png|svg|css|gif|ico|woff2?|ttf)
func ExtensionPass(extension string) bool {
	if extensions.MatchString(extension) {
		return false
	}
	return true
}
