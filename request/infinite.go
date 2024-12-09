package request

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
				Resulted++
				PrintProgress(progressCounter, progressMax)

				requestData := localV.(misc.RequestData)
				iTool.Results(requestData, &m)
				wg.Done()
			}()
		}
		wg.Done()
	}()

	return in
}
