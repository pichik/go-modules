package request

import (
	"regexp"

	"github.com/pichik/go-modules/misc"
	"github.com/pichik/go-modules/tool"
)

var filterFlag, scopeFlag string
var filterRegex, scopeRegex *regexp.Regexp

func (util Filter) SetupFlags() []tool.UtilData {
	var flags []tool.FlagData

	flags = append(flags,
		tool.FlagData{
			Name:        "S",
			Description: "Domain scope regex",
			Required:    false,
			Def:         ".*",
			VarStr:      &scopeFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "F",
			Description: "Unwanted endpoint filtering regex",
			Required:    false,
			Def:         "^!@$",
			VarStr:      &filterFlag,
		})

	examples := make(map[string]string)
	examples["Set domain scope"] = "echo 'google.com' | tt [tool] -S '(admin|support)\\.google.com'"
	examples["Set unwanted endpoit filtering"] = "echo 'google.com' | tt [tool] -F '/(products|blog)/'"

	util.UtilData.Name = "Filters"
	util.UtilData.FlagDatas = flags
	util.UtilData.Examples = examples
	return []tool.UtilData{*util.UtilData}
}

func (util Filter) SetupData() {
	scopeRegex = regexp.MustCompile(scopeFlag)
	filterRegex = regexp.MustCompile(filterFlag)
}

func urlFilterPass(parsedUrl misc.ParsedUrl) bool {
	if scopeRegex != nil && !scopeRegex.MatchString(parsedUrl.Domain) {
		return false
	}
	if filterRegex != nil && filterRegex.MatchString(parsedUrl.Path) {
		return false
	}
	return misc.ExtensionPass(parsedUrl.Extension)
}

// Filter urls from domain scope, uselless paths and extensions
func FilterUrls(urls []misc.ParsedUrl) []misc.ParsedUrl {
	var filteredUrls []misc.ParsedUrl

	for _, url := range urls {
		if urlFilterPass(url) {
			filteredUrls = append(filteredUrls, url)
		}
	}
	return filteredUrls
}
