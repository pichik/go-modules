package request

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/pichik/go-modules/misc"
	"github.com/pichik/go-modules/tool"
)

var outputDirFlag string
var corsCheckFlag, ignoreTestedFlag bool
var threadsFlag int

var sleepTime = int(10)
var currentThreads int
var limiter <-chan time.Time

var progressCounter int
var progressMax int

var slowed bool

func (util misc.RequestFlow) SetupFlags() []tool.UtilData {
	var flags []tool.FlagData

	flags = append(flags,
		tool.FlagData{
			Name:        "t",
			Description: "Threads - requests per second, negative number increase second delay instead of threads (-t -5 :one thread per 5 second)",
			Required:    false,
			Def:         5,
			VarInt:      &threadsFlag,
		})
	flags = append(flags,
		tool.FlagData{
			Name:        "i",
			Description: "Ignore already tested endpoints and clear queue file",
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
	ut = append(ut, Repeater{UtilData: &tool.UtilData{}}.SetupFlags()...)
	ut = append(ut, *util.UtilData)

	return ut
}

func (util RequestFlow) SetupData() {
	var ut []IUtil
	ut = append(ut, Request{})
	ut = append(ut, Repeater{})
	for _, u := range ut {
		u.SetupData()
	}
	//Separate previous found data from current
	misc.DataSeparator()
}

// check what is tested and add rest to queue (dont add it to the tested file)
func SetupQueue(parsedUrls []misc.ParsedUrl, urlSpec string) {
	parsedUrls = FilterUrls(parsedUrls)
	urls := misc.BuildUrls(parsedUrls, urlSpec)
	misc.FillQueue(urls, ignoreTestedFlag)
}

// Simple flow with urls only
func FlowStart(urls []misc.ParsedUrl, iTool IFlowTool) {
	var wg sync.WaitGroup
	in := throughInfinite(&wg, iTool)

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
func CustomFlowStart(requestData []misc.RequestData, iTool IFlowTool) {
	var wg sync.WaitGroup
	in := throughInfinite(&wg, iTool)
	flow(requestData, in)
	close(in)
	wg.Wait()
}

func flow(requestData []misc.RequestData, requestDataChan chan<- interface{}) {
	SetSartTime()

	currentThreads = threadsFlag
	setLimiter(currentThreads)

	var wg sync.WaitGroup
	var m sync.Mutex
	progressCounter = 1
	Resulted = 0
	progressMax = len(requestData) * (1 + Repeats())

	queue := make(chan misc.RequestData, 15)

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, m *sync.Mutex, queue chan misc.RequestData, requestDataChan chan<- interface{}) {
			defer wg.Done()
			for requestData := range queue {
				<-limiter

				//Wait for responses to be resulted, so there is no many unresulted requests
				diff := progressCounter - Resulted
				for diff > 20 {
					time.Sleep(time.Second) // Give up the CPU temporarily
					diff = progressCounter - Resulted
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

	misc.CreateRequest(&requestData)
	if requestData.ResponseStatus == 429 || (requestData.ResponseStatus == 000 && (misc.EOF().MatchString(requestData.Error.Error()) || misc.Timeout().MatchString(requestData.Error.Error()))) {
		misc.PrintUrl(requestData, false)
		slowDown(&check429)
		work(check429, requestData, requestDataChan)
		return
	} else if slowed {
		wg429.Done()
		slowed = false
		if threadsFlag > 1 {
			threadsFlag = threadsFlag / 2
		} else if threadsFlag == 1 {
			threadsFlag = -1
		} else if threadsFlag < 0 {
			threadsFlag--
		}
		currentThreads = threadsFlag
		setLimiter(currentThreads)
		fmt.Printf("Back in game, threads set to: %d\n", currentThreads)
	}

	PrintUrl(requestData, true)
	PrintProgress(progressCounter, progressMax)
	requestDataChan <- requestData
	progressCounter++
	//Send different method if allowed to see if 404 have some valid request methods
	if requestData.Method == "GET" && Repeat(strconv.Itoa(requestData.ResponseStatus)) {
		for _, method := range GetAllMethods() {
			requestData.Method = method
			<-limiter
			work(check429, requestData, requestDataChan)
		}
	}
}

func FlowResults(requestData misc.RequestData, m *sync.Mutex) ([]misc.ScrapData, []misc.ParsedUrl) {
	foundData, completeUrls, incompleteUrls := misc.GetData(requestData.ResponseBody, &requestData.ParsedUrl)

	if requestData.ResponseHeaders.Get("Location") != "" {
		data, comp, incomp := misc.GetData(requestData.ResponseHeaders.Get("Location"), &requestData.ParsedUrl)
		foundData = append(foundData, data...)
		completeUrls = append(completeUrls, comp...)
		incompleteUrls = append(incompleteUrls, incomp...)
	}

	urlToSave := requestData.ParsedUrl.Url
	completeUrls = FilterUrls(completeUrls)
	// m.Lock()
	misc.AddToTested(urlToSave)
	misc.DataOutput(foundData, misc.GetUrls(completeUrls), misc.GetUrls(incompleteUrls))
	misc.CorsOutput(requestData.ResponseHeaders, urlToSave)

	if requestData.ResponseStatus != 404 && requestData.ResponseStatus != 405 && requestData.ResponseContentLength != 0 {
		misc.ResponseOutput(requestData)
	}

	// m.Unlock()

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
		fmt.Printf("Too fast, waiting 30sec to repeat last url\n")
		currentThreads = 1
	} else if *check429 {
		fmt.Printf("Too fast ---Waiting...30sec---\n")
	}
	time.Sleep(time.Duration(30) * time.Second)
}

// create channels that do not require exact number and can go through all urls
func makeInfinite() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})

	go func() {
		var inQueue []interface{}
		outCh := func() chan interface{} {
			if len(inQueue) == 0 {
				return nil
			}
			return out
		}
		curVal := func() interface{} {
			if len(inQueue) == 0 {
				return nil
			}
			return inQueue[0]
		}

		for len(inQueue) > 0 || in != nil {
			select {
			case v, ok := <-in:
				if !ok {
					in = nil
				} else {
					inQueue = append(inQueue, v)
				}
			case outCh() <- curVal():
				inQueue = inQueue[1:]
			}
		}
		close(out)
	}()
	return in, out
}

// go through all urls
func throughInfinite(wg *sync.WaitGroup, iTool IFlowTool) chan<- interface{} {
	var m sync.Mutex

	in, out := makeInfinite()
	wg.Add(1)

	go func() {
		for v := range out {
			wg.Add(1)
			localV := v
			go func() {
				requestData := localV.(misc.RequestData)
				iTool.Results(requestData, &m)
				Resulted++
				PrintProgress(progressCounter, progressMax)
				wg.Done()
			}()
		}
		wg.Done()
	}()

	return in
}
