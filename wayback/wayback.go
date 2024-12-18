package wayback

import (
	"github.com/pichik/go-modules/misc"
	"github.com/pichik/go-modules/request"
	"github.com/pichik/go-modules/tool"
)

type Wayback struct {
	UtilData *tool.UtilData
}

var withTimestampFlag bool
var FilterFlag, fromFlag, pageSizeFlag, ShowDataFlag, SkipByFlag string

func (util Wayback) SetupFlags() []tool.UtilData {
	var flags []tool.FlagData

	flags = append(flags, tool.FlagData{
		Name:        "wt",
		Description: "Get urls with timestamp, use only for specific endpoints, takes too long",
		Required:    false,
		Def:         false,
		VarBool:     &withTimestampFlag,
	})

	flags = append(flags, tool.FlagData{
		Name:        "wfrom",
		Description: "Get data up to this year",
		Required:    false,
		Def:         "2015",
		VarStr:      &fromFlag,
	})
	flags = append(flags, tool.FlagData{
		Name:        "wfilter",
		Description: "Filtering",
		Required:    false,
		Def:         "!mimetype:(font|image)/.*|text/css",
		VarStr:      &FilterFlag,
	})
	flags = append(flags, tool.FlagData{
		Name:        "wsize",
		Description: "Limit response data per request",
		Required:    false,
		Def:         "100",
		VarStr:      &pageSizeFlag,
	})
	flags = append(flags, tool.FlagData{
		Name:        "wdata",
		Description: "What data to show",
		Required:    false,
		Def:         "original,timestamp,mimetype,statuscode,length,digest",
		VarStr:      &ShowDataFlag,
	})
	flags = append(flags, tool.FlagData{
		Name:        "wdigest",
		Description: "What data to show",
		Required:    false,
		Def:         "digest",
		VarStr:      &SkipByFlag,
	})

	util.UtilData.Name = "HTTP Request"
	util.UtilData.FlagDatas = flags

	var ut []tool.UtilData
	ut = append(ut, request.RequestFlow{UtilData: &tool.UtilData{}}.SetupFlags()...)
	ut = append(ut, *util.UtilData)

	return ut
}

func (util Wayback) SetupData() {
	request.NormalizeFlag = false
	request.ThreadsFlag = -4

	var ut []tool.IUtil
	ut = append(ut, request.RequestFlow{})
	for _, u := range ut {
		u.SetupData()
	}

}

func HandleResponse(requestData misc.RequestData) ([]string, *misc.ParsedUrl) {

	// check if there's results, used instead of pagination, as wayback's pagination response is not always correct when using a filter
	if requestData.ResponseStatus != 200 {
		return []string{}, nil
	}

	res := []WB{}
	UnmarshalWB(requestData.ResponseBodyBytes, &res)

	rawUrls := getRawUrls(res)

	//Prevent creating additional requests if response contains only few urls
	if requestData.ResponseContentLength < 10000 {
		return rawUrls, nil
	}
	nextUrl := buildNextPageUrl(requestData.ParsedUrl)

	return rawUrls, nextUrl
}
