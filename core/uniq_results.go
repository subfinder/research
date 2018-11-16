package core

// UniqResults filters a given input stream for uniq outputs. Note: this
// will only be a filter for successful results.
func UniqResults(input <-chan *Result) <-chan *Result {
	output := make(chan *Result)

	go func() {
		defer close(output)

		uniqFilter := map[string]bool{}

		for i := range input {
			if i.IsSuccess() {
				str, ok := i.Success.(string)
				if !ok {
					continue
				}
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					output <- i
				}
			}
		}
	}()

	return output
}
