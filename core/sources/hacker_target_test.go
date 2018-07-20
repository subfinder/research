package sources

import core "github.com/subfinder/research/core"
import "testing"
import "fmt"
import "sync"
import "strings"

func TestHackerTarget(t *testing.T) {
	domain := "google.com"
	source := HackerTarget{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if len(results) == 1 {
		if !strings.Contains(results[0].Failure.Error(), "API count exceeded") {
			t.Errorf("expected to return API count error, got '%v'", results[0].Failure.Error())
		}
	} else {
		if !(len(results) >= 4000) {
			t.Errorf("expected to return more than one successful result, got %v", len(results))
		}
	}
}

func TestHackerTarget_multi_threaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := HackerTarget{}
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

	if len(results) == 4 {
		if !strings.Contains(results[0].Failure.Error(), "API count exceeded") {
			t.Errorf("expected to return API count error, got '%v'", results[0].Failure.Error())
		}
	} else {
		if !(len(results) >= 4000) {
			t.Errorf("expected to return more than one successful result, got %v", len(results))
		}
	}
}

func ExampleHackerTarget() {
	domain := "google.com"
	source := HackerTarget{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	fmt.Println(len(results) >= 1)
	// Output: true
}

func ExampleHackerTarget_multi_threaded() {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := HackerTarget{}
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

	fmt.Println(len(results) >= 1)
	// Output: true
}

func BenchmarkHackerTargetSingleThreaded(b *testing.B) {
	domain := "google.com"
	source := HackerTarget{}

	for n := 0; n < b.N; n++ {
		results := []*core.Result{}
		for result := range source.ProcessDomain(domain) {
			results = append(results, result)
		}
	}
}

func BenchmarkHackerTargetMultiThreaded(b *testing.B) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := HackerTarget{}
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	for n := 0; n < b.N; n++ {
		results := []*core.Result{}

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
	}
}
