package core

import (
	"context"
	"fmt"
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
	results := make(chan *Result)
	go func() {
		defer fmt.Println("*** done sending results ***")
		defer close(results)
		wg := sync.WaitGroup{}
		for _, source := range options.Sources {
			wg.Add(1)
			go func(source Source) {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				sourceResults := source.ProcessDomain(domain)
				for {
					select {
					case <-time.After(15 * time.Second):
						fmt.Println("*** time after ***")
						return
					case result, ok := <-sourceResults:
						if ok {
							select {
							case results <- result:
								// no timeout
							case <-time.After(15 * time.Second):
								fmt.Println("*** time after on pass along ***")
								return
							}
							results <- result
						} else {
							fmt.Println("*** not ok ***")
							return
						}
					case <-ctx.Done():
						fmt.Println("*** ctx done ***")
						return
					}
				}
			}(source)
		}
		wg.Wait()
	}()
	return results
}
