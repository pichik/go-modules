package request

import (
	"strings"

	"github.com/pichik/go-modules/tool"
)

var repeatsFlag string
var methodsRepeat []string
var statusRepeat []string
var repeats int

func (util Repeater) SetupFlags() []tool.UtilData {
	var flags []tool.FlagData

	flags = append(flags,
		tool.FlagData{
			Name:        "rep",
			Description: "Repeat request with different methods",
			Required:    false,
			Def:         "",
			VarStr:      &repeatsFlag,
		})

	examples := make(map[string]string)
	examples["Repeat all methods to every response status"] = "secho 'google.com' | tt crawl -rep all-all"
	examples["Repeat specific methods to a specific response status"] = "echo 'google.com' | tt crawl -rep 404,405-put,post"

	util.UtilData.Name = "Repeater"
	util.UtilData.FlagDatas = flags
	util.UtilData.Examples = examples
	return []tool.UtilData{*util.UtilData}
}

func (util Repeater) SetupData() {
	if repeatsFlag == "" {
		return
	}
	statusMethod := strings.Split(repeatsFlag, "-")

	if len(statusMethod) < 2 {
		return
	}
	statusRepeat = strings.Split(statusMethod[0], ",")
	methodsRepeat = strings.Split(statusMethod[1], ",")
	checkMethods()

	if statusRepeat[0] == "all" {
		repeats = len(methodsRepeat)
	}
}

func GetAllMethods() []string {
	return methodsRepeat
}

func SetupRepeaterData() {
	if repeatsFlag == "" {
		return
	}
	statusMethod := strings.Split(repeatsFlag, "-")

	if len(statusMethod) < 2 {
		return
	}
	statusRepeat = strings.Split(statusMethod[0], ",")
	methodsRepeat = strings.Split(statusMethod[1], ",")
	checkMethods()

	if statusRepeat[0] == "all" {
		repeats = len(methodsRepeat)
	}

}

func Repeat(currentStatus string) bool {

	if len(statusRepeat) == 0 {
		return false
	}
	if statusRepeat[0] == "all" {
		return true
	}
	for _, status := range statusRepeat {
		if status == currentStatus {
			return true
		}
	}
	return false
}

func Repeats() int {
	return repeats
}

func checkMethods() {
	if methodsRepeat[0] == "all" {
		methodsRepeat[0] = "POST"
		methodsRepeat = append(methodsRepeat, "PUT")
		methodsRepeat = append(methodsRepeat, "PATCH")
		methodsRepeat = append(methodsRepeat, "DELETE")
	}
}
