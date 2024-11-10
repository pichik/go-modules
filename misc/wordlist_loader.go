package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var ApiData API

type API struct {
	Objects []string
	Actions []string
}

var FilterData []Filters

type Filters struct {
	Name        string
	RegexString string
	RegexPart   int
	Highlight   string
}

var EndpointData = make(map[string][]string)

var Bypasses []string

var wordlistUrl = string("raw.githubusercontent.com/pichik/wordlists/main/")

func LoadFilters() {
	js, _ := loadWordlist("filters.json")
	err := json.Unmarshal(js, &FilterData)

	errorCheck("filters", err)
}

func LoadApis() {
	_, ApiData.Actions = loadWordlist("apis/actions.txt")
	_, ApiData.Objects = loadWordlist("apis/objects.txt")
}

func LoadBypasses() {
	_, Bypasses = loadWordlist("bypasses.txt")
}

func LoadEndpoints(tag string) {
	_, EndpointData[tag] = loadWordlist(fmt.Sprintf("endpoints/%s.txt", tag))
}

func loadWordlist(file string) ([]byte, []string) {
	url := fmt.Sprintf("https://%s%s", wordlistUrl, file)
	res, err := http.Get(url)

	errorCheck(file, err)

	if res.StatusCode == 404 {
		err = errors.New("wordlist not found")
	}

	errorCheck(file, err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	errorCheck(file, err)

	//create list from every line
	list := strings.Fields(string(body))

	return body, list
}

// Use for debugging only !!!
// func loadLocal(file string) []string {
// 	fmt.Println("!!! LOCAL WORDLIST - DEBUG ONLY !!!")

// 	result, err := Read(file)

// 	errorCheck(file, err)

// 	return result
// }

func errorCheck(wordlist string, err error) {
	if err == nil {
		return
	}
	fmt.Printf("%sUnable to retreive %s:\n\t%s%s\nCheck %shttps://github.com/pichik/thetool/tree/main/wordlists%s for available wordlists\n", Red, wordlist, err, White, Blue, White)
	os.Exit(0)
}
