package request

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/pichik/go-modules/misc"
	"github.com/pichik/go-modules/parser"
	"github.com/pichik/go-modules/request/data"
	"github.com/pichik/go-modules/tool"
)

type RequestFlow struct {
	UtilData *tool.UtilData
}
type IFlowTool interface {
	Results(requestData misc.RequestData, m *sync.Mutex)
}

var outputDirFlag string
var corsCheckFlag, ignoreTestedFlag bool
var ThreadsFlag int

var sleepTime = int(10)
var limiter <-chan time.Time

// var progressCounter int
// var progressMax int

var slowed bool

var saveResults bool

func (util RequestFlow) SetupFlags() []tool.UtilData {
	var flags []tool.FlagData

	flags = append(flags,
		tool.FlagData{
			Name:        "t",
			Description: "Threads - requests per second, negative number increase second delay instead of threads (-t -5 :one thread per 5 second)",
			Required:    false,
			Def:         5,
			VarInt:      &ThreadsFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "i",
			Description: "Ignore queue and tested urls",
			Required:    false,
			Def:         false,
			VarBool:     &ignoreTestedFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "o",
			Description: fmt.Sprintf("Output directory \n\t\t%s(Default: tool name)%s", misc.Gray, misc.White),
			Required:    false,
			Def:         "",
			VarStr:      &outputDirFlag,
		})

	util.UtilData.FlagDatas = flags
	var ut []tool.UtilData
	ut = append(ut, Request{UtilData: &tool.UtilData{}}.SetupFlags()...)
	ut = append(ut, parser.Parser{UtilData: &tool.UtilData{}}.SetupFlags()...)
	ut = append(ut, Repeater{UtilData: &tool.UtilData{}}.SetupFlags()...)
	ut = append(ut, *util.UtilData)

	return ut
}

func (util RequestFlow) SetupData() {
	var ut []tool.IUtil
	ut = append(ut, Request{})
	ut = append(ut, parser.Parser{})
	ut = append(ut, Repeater{})
	for _, u := range ut {
		u.SetupData()
	}

	if ignoreTestedFlag {
		misc.RemoveQueueFile()
	}
}

// check what is tested and add rest to queue (dont add it to the tested file)
func SetupQueue(parsedUrls []misc.ParsedUrl, urlSpec string) {
	parsedUrls = parser.FilterUrls(parsedUrls)
	urls := misc.BuildUrls(parsedUrls, urlSpec)
	misc.FillQueue(urls, ignoreTestedFlag)
}

// Simple flow with urls only
func FlowStart(urls []misc.ParsedUrl, iTool IFlowTool, saveRes bool) {
	saveResults = saveRes
	var wg sync.WaitGroup
	// in := throughInfinite(&wg, iTool)
	in := data.ThroughInfinite(&wg, func(requestData misc.RequestData, m *sync.Mutex) {
		iTool.Results(requestData, m) // Use your tool here
	})

	var requestData = make([]misc.RequestData, len(urls))
	for i, u := range urls {
		requestData[i] = RequestBase
		requestData[i].ParsedUrl = u
		//If user didnt set origin, set it to requested domainz
		if requestData[i].Headers["Origin"] == "" {
			requestData[i].Headers["Origin"] = misc.BuildUrl(requestData[i].ParsedUrl, "12")
		}
	}

	flow(requestData, in)
	close(in)
	wg.Wait()
}

// Flow with customable request data, headers etc..
func CustomFlowStart(requestData []misc.RequestData, iTool IFlowTool, save bool) {
	saveResults = save
	var wg sync.WaitGroup
	// in := throughInfinite(&wg, iTool)
	in := data.ThroughInfinite(&wg, func(requestData misc.RequestData, m *sync.Mutex) {
		iTool.Results(requestData, m) // Use your tool here
	})
	flow(requestData, in)
	close(in)
	wg.Wait()
}

func flow(requestData []misc.RequestData, requestDataChan chan<- interface{}) {
	data.SetSartTime()

	data.CurrentThreads = ThreadsFlag
	setLimiter(data.CurrentThreads)

	var wg sync.WaitGroup
	var m sync.Mutex
	data.ProgressCounter = 0
	data.Resulted = 0
	data.ProgressMax = len(requestData) * (1 + Repeats())

	queue := make(chan misc.RequestData, 15)

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, m *sync.Mutex, queue chan misc.RequestData, requestDataChan chan<- interface{}) {
			defer wg.Done()
			for requestData := range queue {
				<-limiter

				//Wait for responses to be resulted, so there is no many unresulted requests
				diff := data.ProgressCounter - data.Resulted
				for diff > 20 {
					time.Sleep(time.Second) // Give up the CPU temporarily
					diff = data.ProgressCounter - data.Resulted
				}

				work(false, requestData, requestDataChan)
			}
		}(&wg, &m, queue, requestDataChan)
	}

	for _, rd := range requestData {
		queue <- rd
	}

	close(queue)
	wg.Wait()
}

var wg429 sync.WaitGroup

func work(check429 bool, requestData misc.RequestData, requestDataChan chan<- interface{}) {

	if slowed && !check429 {
		wg429.Wait()
	}

	CreateRequest(&requestData)
	if requestData.ResponseStatus == 429 || (requestData.ResponseStatus == 000 && misc.RepeatRequestTriggers().MatchString(requestData.Error.Error())) {
		data.PrintUrl(requestData, false)
		slowDown(&check429)
		work(check429, requestData, requestDataChan)
		return
	} else if slowed {
		wg429.Done()
		slowed = false
		if ThreadsFlag > 1 {
			ThreadsFlag = ThreadsFlag / 2
		} else if ThreadsFlag == 1 {
			ThreadsFlag = -1
		} else if ThreadsFlag < 0 {
			ThreadsFlag--
		}
		data.CurrentThreads = ThreadsFlag
		setLimiter(data.CurrentThreads)
	}

	data.ProgressCounter++
	data.PrintUrl(requestData, saveResults)
	data.PrintProgress()
	requestDataChan <- requestData
	//Send different method if allowed to see if 404 have some valid request methods
	if requestData.Method == "GET" && Repeat(strconv.Itoa(requestData.ResponseStatus)) {
		for _, method := range GetAllMethods() {
			requestData.Method = method
			<-limiter
			work(check429, requestData, requestDataChan)
		}
	}
}

func FlowResults(requestData misc.RequestData, m *sync.Mutex) (map[string]parser.ParserData, []misc.ParsedUrl) {

	foundData, completeUrls := parser.ParseText(requestData.ResponseBody, &requestData.ParsedUrl)

	if requestData.ResponseHeaders.Get("Location") != "" {
		data, comp := parser.ParseText(requestData.ResponseHeaders.Get("Location"), &requestData.ParsedUrl)
		foundData = parser.MergeData(foundData, data)
		completeUrls = append(completeUrls, comp...)
	}

	return foundData, completeUrls
}

func setLimiter(threads int) {
	if threads > 0 {
		limiter = time.Tick(time.Duration(1000/threads) * time.Millisecond)
	} else if threads < 0 {
		limiter = time.Tick(time.Duration(threads*-1000) * time.Millisecond)
	}
}

func slowDown(check429 *bool) {
	if !slowed {
		*check429 = true
		wg429.Add(1)
		slowed = true
		data.CurrentThreads = 1
	}
	time.Sleep(time.Duration(30) * time.Second)
}
