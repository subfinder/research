package core

import (
	"context"
	"sync"
	"time"
)

// EnumerateSubdomains takes the given domain and with each Source from EnumerationOptions,
// it will spawn a go routine to start processing that Domain. The result channels from each
// source are merged into one results channel to be consumed.
//
//
//
//   ____________________________     Source1.ProcessDomain     ___________        _____
//  |                            | /                         \ |           |      |     |
//  | EnumerationOptions.Sources | -- Source2.ProcessDomain -- |  Results  | ---> |  ?  |
//  |____________________________| \                         / |___________|      |_____|
//                                    Source3.ProcessDomain
//
//
func EnumerateSubdomains(domain string, options *EnumerationOptions) <-chan *Result {
	// this channel of results will be used to combine the result channels
	// from each source configured in the EnumerationOptions
	results := make(chan *Result)
	// the main processing will be done in the background
	go func() {
		// close up the combined results channel when the go parent func returns
		defer close(results)

		// a wait group to ensure all child go funcs finish processing
		wg := sync.WaitGroup{}

		// iterate over each source provided in the EnumerationOptions
		for _, source := range options.Sources {

			// register a job in the wait group
			wg.Add(1)

			// spawn a go func with the current source in the iteration
			go func(source Source) {

				// Tell the wait group a job has been completed when the go func returns
				defer wg.Done()

				// ctx is used for ensuring there are no lingering go funcs to avoid memory leaks
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				// get the results channel from the source calling the ProcessDomain method on it
				sourceResults := source.ProcessDomain(domain)

				// for loop over results in a select to allow for timeout
				for {
					select {
					case result, ok := <-sourceResults:
						if ok {
							select {
							case results <- result:
								// no timeout
							case <-ctx.Done():
								// timed out while passing result to combined results channel
								return
							}
						}
					case <-ctx.Done():
						// timed out while getting a result from the source's results channel
						return
					}
				}
			}(source)
		}
		wg.Wait()
	}()
	// this function returns the combined results channel right away
	return results
}
