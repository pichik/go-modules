package wayback

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pichik/go-modules/misc"
)

var timestampUrlBase = "https://web.archive.org/web/"
var archiveUrlBase = "https://web.archive.org/cdx/search/cdx?output=json"

var timestampRegex = regexp.MustCompile(`/web/(\d{14})if_/`)
var httpRegex = regexp.MustCompile(`if_/(http.*)`)

func BuildUrls(urls *[]string) {
	archiveUrlBase = fmt.Sprintf("%s&pageSize=%s&filter=%s&from=%s&fl=%s&collapse=%s&page=0", archiveUrlBase, pageSizeFlag, FilterFlag, fromFlag, ShowDataFlag, SkipByFlag)

	for i, url := range *urls {
		(*urls)[i] = fmt.Sprintf("%s&url=%s*", archiveUrlBase, url)
	}
}

func BuildTimestampUrls(wb []WB) []string {
	var tsUrls []string

	for _, entry := range wb {
		timestampUrl := fmt.Sprintf("%s%sif_/%s", timestampUrlBase, entry.Timestamp, entry.Original)
		tsUrls = append(tsUrls, timestampUrl)

	}
	return tsUrls
}

func RemoveArchiveUrl(archiveUrl misc.ParsedUrl) (string, misc.ParsedUrl) {
	archiveEndpoint := misc.BuildUrl(archiveUrl, "3")
	timestamp := timestampRegex.FindStringSubmatch(archiveEndpoint)[1]
	url := httpRegex.FindStringSubmatch(archiveEndpoint)[1]
	return timestamp, misc.ParseUrl(url)
}

func getRawUrls(wbs []WB) []string {
	var rawUrls []string

	for _, wb := range wbs {
		url := strings.TrimSuffix(wb.Original, "/")
		if misc.ExtensionPass(url) {
			//Parse url to remove :80 protocols to be appended
			purl := misc.ParseUrl(url)
			purl.Protocol = "https"
			misc.RebuildUrl(&purl)
			rawUrls = append(rawUrls, purl.Url)
		}
	}
	return rawUrls
}

// func buildResumeUrl(parsedUrl misc.ParsedUrl, resumeKey string) *misc.ParsedUrl {
// 	if resumeKey == "" {
// 		return nil
// 	}

// 	parsedUrl.Queries.Set("resumeKey", resumeKey)
// 	misc.RebuildUrl(&parsedUrl)
// 	return &parsedUrl
// }

func buildNextPageUrl(parsedUrl misc.ParsedUrl) *misc.ParsedUrl {
	currPage := parsedUrl.Queries.Get("page")

	page, err := strconv.Atoi(currPage)
	if err != nil {
		misc.PrintError("Pagination to Int", err)
		return nil
	}
	page++

	parsedUrl.Queries.Set("page", strconv.Itoa(page))
	misc.RebuildUrl(&parsedUrl)
	return &parsedUrl
}
