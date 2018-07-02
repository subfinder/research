package subzero

func AggregateSuccessfulResults(in chan *Result) <-chan *Result {
	out := make(chan *Result)
	go func(in, out chan *Result) {
		defer close(out)
		for result := range in {
			if result.Failure == nil {
				out <- result
			}
		}
	}(in, out)
	return out
}

