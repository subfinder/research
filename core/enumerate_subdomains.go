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
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		for _, source := range options.Sources {
			wg.Add(1)
			go func(source Source, ctx context.Context, cancel context.CancelFunc) {
				defer wg.Done()
				sourceResults := source.ProcessDomain(domain)
				for {
					select {
					case result, ok := <-sourceResults:
						if ok {
							results <- result
						} else {
							return
						}
					case <-ctx.Done():
						return
					}
				}
			}(source, ctx, cancel)
		}
		wg.Wait()
	}()
	return results
}
