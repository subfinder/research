package core

import (
	"context"
	"sync"
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
func EnumerateSubdomains(ctx context.Context, domain string, options *EnumerationOptions) <-chan *Result {
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

				// get the results channel from the source calling the ProcessDomain method on it
				sourceResults := source.ProcessDomain(ctx, domain)

				// for loop over results in a select to allow for timeout
				for {
					select {
					case result, ok := <-sourceResults:
						if ok {
							select {
							case results <- result:
								// initial recursion implementation
								if options.Recursive && result.IsSuccess() {
									str, ok := result.Success.(string)
									if !ok {
										continue
									}
									wg.Add(1)

									go func(results chan *Result, domain string, options *EnumerationOptions) {
										defer wg.Done()
										for result := range EnumerateSubdomains(ctx, domain, options) {
											select {
											case <-ctx.Done():
												return
											case results <- result:
												continue
											}
										}
									}(results, str, options)
								}
								// no timeout
							case <-ctx.Done():
								// timed out while passing result to combined results channel
								return
							}
						} else {
							// failed to retrieve result from results channel
							return
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

	if options.Uniq {
		return UniqResults(results)
	}

	// this function returns the combined results channel right away
	return results
}
