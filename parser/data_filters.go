package parser

import (
	"sync"

	"github.com/pichik/go-modules/misc"
)

func getData(text string, parserData *ParserData) {
	var wg sync.WaitGroup // To wait for all goroutines to finish
	mu := sync.Mutex{}    // Mutex to protect the shared parserData.Results
	for _, filter := range parserData.Filters {
		wg.Add(1) // Add a goroutine to the wait group

		go func(filter Filters) {
			defer wg.Done() // Mark this goroutine as done when it finishes

			// Run ag with the current filter's RegexString
			findings, err := runAg(text, filter.RegexString)
			if err != nil {
				misc.PrintError("Regex failed", err)
				return
			}

			// Iterate over the findings and highlight the matches in the goroutine
			for i, finding := range findings {
				findings[i] = misc.Highlight(finding, filter.Highlight)
			}

			// Lock the shared parserData.Results before modifying it
			mu.Lock()
			parserData.Results = append(parserData.Results, findings...)
			mu.Unlock()

		}(filter) // Pass the filter as an argument to the goroutine
	}

	// Wait for all goroutines to finish before returning
	wg.Wait()
}
