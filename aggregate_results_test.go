package subzero

import "testing"
import "fmt"
import "strings"
import "errors"

func TestAggregateSuccessfulResults(t *testing.T) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: 0},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	counter := 0

	for _ = range AggregateSuccessfulResults(fakeResultsChan) {
		counter++
	}

	if counter != 3 {
		t.Fatalf("expected '%v' successful results, got '%v'", 3, counter)
	}
}

func TestAggregateFailedResults(t *testing.T) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: 0},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	counter := 0

	for _ = range AggregateFailedResults(fakeResultsChan) {
		counter++
	}

	if counter != 2 {
		t.Fatalf("expected '%v' failed results, got '%v'", 2, counter)
	}
}

func TestAggregateCustomResults(t *testing.T) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: false},
		&Result{Success: 0},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	counter := 0

	for _ = range AggregateCustomResults(fakeResultsChan, func(r *Result) bool {
		_, ok := r.Success.(bool)
		return ok
	}) {
		counter++
	}

	if counter != 2 {
		t.Fatalf("expected '%v' successful results, got '%v'", 2, counter)
	}
}

func TestAggregateCustomResultsMore(t *testing.T) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: false},
		&Result{Success: 0},
		&Result{Success: "picat"},
		&Result{Success: "was"},
		&Result{Success: "here"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	puzzle := []string{}

	for result := range AggregateCustomResults(fakeResultsChan, func(r *Result) bool {
		_, ok := r.Success.(string)
		return ok
	}) {
		puzzle = append(puzzle, result.Success.(string))
	}

	if strings.Join(puzzle, " ") != "picat was here" {
		t.Fatalf("expected '%v', got '%v'", "picat was here", strings.Join(puzzle, " "))
	}
}

func ExampleAggregateCustomResults() {
	fakeResults := []*Result{
		&Result{Type: "color", Success: "red"},
		&Result{Type: "color", Success: "green"},
		&Result{Type: "color", Success: "blue"},
		&Result{Type: "color", Failure: errors.New("no color")},
		&Result{Type: "wiggle", Failure: errors.New("wiggle")},
		&Result{Success: "example"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	// imagine a function doing something useful
	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	// consume aggregated results
	for result := range AggregateCustomResults(fakeResultsChan, func(r *Result) bool {
		return r.Type == "color" && r.IsSuccess() // only successful results of type "color"
	}) {
		fmt.Println(result.Success)
	}
	// Output:
	// red
	// green
	// blue
}

func ExampleAggregateSuccessfulResults() {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: 0},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	counter := 0

	for _ = range AggregateSuccessfulResults(fakeResultsChan) {
		counter++
	}

	fmt.Println(counter)
	// Output: 3
}

func ExampleAggregateFailedResults() {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: 0},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	counter := 0

	for _ = range AggregateFailedResults(fakeResultsChan) {
		counter++
	}

	fmt.Println(counter)
	// Output: 2
}

func BenchmarkAggregateSuccessfulResults(b *testing.B) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Failure: errors.New("example1")},
		&Result{Success: 0},
		&Result{Failure: errors.New("example2")},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example3")},
	}
	for n := 0; n < b.N; n++ {
		fakeResultsChan := make(chan *Result)
		go func(fakeResults []*Result, fakeResultsChan chan *Result) {
			defer close(fakeResultsChan)
			for _, result := range fakeResults {
				fakeResultsChan <- result
			}
		}(fakeResults, fakeResultsChan)

		counter := 0

		for _ = range AggregateSuccessfulResults(fakeResultsChan) {
			counter++
		}
	}
}

func BenchmarkAggregateFailedResults(b *testing.B) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Failure: errors.New("example1")},
		&Result{Success: 0},
		&Result{Failure: errors.New("example2")},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example3")},
	}

	for n := 0; n < b.N; n++ {
		fakeResultsChan := make(chan *Result)
		go func(fakeResults []*Result, fakeResultsChan chan *Result) {
			defer close(fakeResultsChan)
			for _, result := range fakeResults {
				fakeResultsChan <- result
			}
		}(fakeResults, fakeResultsChan)

		counter := 0

		for _ = range AggregateFailedResults(fakeResultsChan) {
			counter++
		}
	}
}

func BenchmarkAggregateCustomResultsSuccessful(b *testing.B) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Failure: errors.New("example1")},
		&Result{Success: 0},
		&Result{Failure: errors.New("example2")},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example3")},
	}

	var successfulOnly = func(r *Result) bool {
		return r.IsSuccess()
	}

	for n := 0; n < b.N; n++ {
		fakeResultsChan := make(chan *Result)
		go func(fakeResults []*Result, fakeResultsChan chan *Result) {
			defer close(fakeResultsChan)
			for _, result := range fakeResults {
				fakeResultsChan <- result
			}
		}(fakeResults, fakeResultsChan)

		counter := 0

		for _ = range AggregateCustomResults(fakeResultsChan, successfulOnly) {
			counter++
		}
	}
}

func BenchmarkAggregateCustomResultsFailed(b *testing.B) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Failure: errors.New("example1")},
		&Result{Success: 0},
		&Result{Failure: errors.New("example2")},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example3")},
	}

	var successfulOnly = func(r *Result) bool {
		return r.IsFailure()
	}

	for n := 0; n < b.N; n++ {
		fakeResultsChan := make(chan *Result)
		go func(fakeResults []*Result, fakeResultsChan chan *Result) {
			defer close(fakeResultsChan)
			for _, result := range fakeResults {
				fakeResultsChan <- result
			}
		}(fakeResults, fakeResultsChan)

		counter := 0

		for _ = range AggregateCustomResults(fakeResultsChan, successfulOnly) {
			counter++
		}
	}
}

func BenchmarkAggregateCustomResultsSuccessfulStrings(b *testing.B) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Failure: errors.New("example1")},
		&Result{Success: 0},
		&Result{Failure: errors.New("example2")},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example3")},
	}

	var successfulStringsOnly = func(r *Result) bool {
		_, ok := r.Success.(string)
		return ok
	}

	for n := 0; n < b.N; n++ {
		fakeResultsChan := make(chan *Result)
		go func(fakeResults []*Result, fakeResultsChan chan *Result) {
			defer close(fakeResultsChan)
			for _, result := range fakeResults {
				fakeResultsChan <- result
			}
		}(fakeResults, fakeResultsChan)

		counter := 0

		for _ = range AggregateCustomResults(fakeResultsChan, successfulStringsOnly) {
			counter++
		}
	}
}

