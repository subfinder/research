package subzero

import "sync"

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
	results := make(chan *Result)
	go func(domain string, options *EnumerationOptions, results chan *Result) {
		defer close(results)
		wg := sync.WaitGroup{}
		for _, source := range options.Sources {
			wg.Add(1)
			// merge source results channels with go funcs
			go func(source Source, results chan *Result) {
				defer wg.Done()
				for result := range source.ProcessDomain(domain) {
					results <- result
				}
			}(source, results)
		}
		wg.Wait()
	}(domain, options, results)
	return results
}
