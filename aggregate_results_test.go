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

