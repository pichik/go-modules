package misc

import (
	"regexp"
	"strings"
)

var extensions = regexp.MustCompile(`(pdf|jpe?g|png|svg|css|gif|ico|woff2?|ttf)`)
var repeatRequestTriggers = regexp.MustCompile(`(context deadline exceeded|EOF|dial tcp.*i/o timeout|connection refused)$`)

func RepeatRequestTriggers() *regexp.Regexp {
	return repeatRequestTriggers
}

// Is not uselless file extension
func ExtensionPass(protocol string, extension string) bool {
	if !strings.Contains(protocol, "http") || extensions.MatchString(extension) {
		return false
	}
	return true
}

// check if domain is missing in endpoints and add domain
func AddDomain(domain string, parsedUrl ...*ParsedUrl) {
	if domain == "" {
		return
	}
	for _, pu := range parsedUrl {
		if pu.Domain == "" {
			pu.Domain = domain
			RebuildUrl(pu)
		}
	}
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
