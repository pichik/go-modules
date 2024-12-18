package data

import (
	"sync"

	"github.com/pichik/go-modules/misc"
)

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

type ProcessFunc func(misc.RequestData, *sync.Mutex)

func ThroughInfinite(wg *sync.WaitGroup, process ProcessFunc) chan<- interface{} {
	var m sync.Mutex

	in, out := makeInfinite()
	wg.Add(1)

	go func() {
		for v := range out {
			wg.Add(1)
			localV := v
			go func() {
				Resulted++
				PrintProgress()

				requestData := localV.(misc.RequestData)
				process(requestData, &m) // Call the provided function
				wg.Done()
			}()
		}
		wg.Done()
	}()

	return in
}

// // go through all urls
// func throughInfinite(wg *sync.WaitGroup, iTool IFlowTool) chan<- interface{} {
// 	var m sync.Mutex

// 	in, out := makeInfinite()
// 	wg.Add(1)

// 	go func() {
// 		for v := range out {
// 			wg.Add(1)
// 			localV := v
// 			go func() {
// 				data.Resulted++
// 				data.PrintProgress(progressCounter, progressMax)

// 				requestData := localV.(misc.RequestData)
// 				iTool.Results(requestData, &m)
// 				wg.Done()
// 			}()
// 		}
// 		wg.Done()
// 	}()

// 	return in
// }
