package wayback

import (
	"encoding/json"
	"fmt"

	"github.com/pichik/go-modules/misc"
)

type WB struct {
	Timestamp  string
	Original   string
	Mimetype   string
	Statuscode string
	Length     string
	Digest     string
}

var collectedWBData []WB

func UnmarshalWB(body []byte, res *[]WB) {
	var records [][]interface{} // To handle the array of arrays structure
	if err := json.Unmarshal(body, &records); err != nil {
		misc.PrintError("Unmarshal wb", err)
		return
	}
	if len(records) < 1 {
		return
	}

	var results []WB
	for _, record := range records[1:] { // Skip the header row (first row)

		wb := WB{}
		// Assign fields based on available length in `record`
		if len(record) > 0 {
			wb.Original = fmt.Sprintf("%v", record[0]) // Original URL
		}
		if len(record) > 1 {
			wb.Timestamp = fmt.Sprintf("%v", record[1]) // Timestamp
		}
		if len(record) > 2 {
			wb.Mimetype = fmt.Sprintf("%v", record[2]) // Mimetype
		}
		if len(record) > 3 {
			wb.Statuscode = fmt.Sprintf("%v", record[3]) // Statuscode
		}
		if len(record) > 4 {
			wb.Length = fmt.Sprintf("%v", record[4]) // Length
		}
		if len(record) > 5 {
			wb.Digest = fmt.Sprintf("%v", record[5]) // Length
		}

		results = append(results, wb)
		collectedWBData = append(collectedWBData, wb)

	}

	*res = results
}

func GetCollectedWBData() []WB {
	return collectedWBData
}
