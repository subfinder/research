package core

import (
	"sync"
)

// MergeResults takes in N number of result channels and merges them
// into one resul channel.
func MergeResults(inputs ...<-chan *Result) <-chan *Result {
	output := make(chan *Result)

	go func() {
		defer close(output)

		wg := sync.WaitGroup{}

		wg.Add(len(inputs))

		for _, input := range inputs {
			go func(input <-chan *Result) {
				defer wg.Done()
				for i := range input {
					output <- i
				}
			}(input)
		}

		wg.Wait()
	}()

	return output
}
