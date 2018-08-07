package sources

import (
	"fmt"
	"sync"
	"testing"

	"github.com/subfinder/research/core"
)

func TestBing(t *testing.T) {
	domain := "google.com"
	source := Bing{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		t.Log(result)
		results = append(results, result)
		// Not waiting around to iterate all the possible pages.
		if len(results) >= 20 {
			break
		}
	}

	if !(len(results) >= 20) {
		t.Errorf("expected more than 20 result(s), got '%v'", len(results))
	}
}

func TestBing_multi_threaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := Bing{}
	results := []*core.Result{}

	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			for result := range source.ProcessDomain(domain) {
				t.Log(result)
				if result.IsSuccess() && result.IsFailure() {
					t.Error("got a result that was a success and failure")
				}
				mx.Lock()
				results = append(results, result)
				mx.Unlock()
			}
		}(domain)
	}

	wg.Wait() // collect results

	if len(results) <= 4 {
		t.Errorf("expected at least 4 results, got '%v'", len(results))
	}
}

func ExampleBing() {
	domain := "google.com"
	source := Bing{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	fmt.Println(len(results) >= 20)
	// Output: true
}

func ExampleBing_multi_threaded() {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := Bing{}
	results := []*core.Result{}

	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			for result := range source.ProcessDomain(domain) {
				mx.Lock()
				results = append(results, result)
				mx.Unlock()
			}
		}(domain)
	}

	wg.Wait() // collect results

	fmt.Println(len(results) >= 4)
	// Output: true
}

