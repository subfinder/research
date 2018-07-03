package subzero

// AggregateSuccessfulResults takes a given results channel as input
// and only sends successful results down the returned output channel.
func AggregateSuccessfulResults(in chan *Result) <-chan *Result {
	out := make(chan *Result)
	go func(in, out chan *Result) {
		defer close(out)
		for result := range in {
			if result.IsSuccess() {
				out <- result
			}
		}
	}(in, out)
	return out
}

// AggregateFailedResults takes a given results channel as input
// and only sends failed results down the returned output channel.
func AggregateFailedResults(in chan *Result) <-chan *Result {
	out := make(chan *Result)
	go func(in, out chan *Result) {
		defer close(out)
		for result := range in {
			if result.Failure != nil {
				out <- result
			}
		}
	}(in, out)
	return out
}

// AggregateCustomResults takes a given results channel as input
// along with a custom filter function that will be executed with each Result.
func AggregateCustomResults(in chan *Result, custom func(r *Result) bool) <-chan *Result {
	out := make(chan *Result)
	go func(in, out chan *Result) {
		defer close(out)
		for result := range in {
			if custom(result) {
				out <- result
			}
		}
	}(in, out)
	return out
}
