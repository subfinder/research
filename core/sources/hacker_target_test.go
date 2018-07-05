package sources

import "testing"
import "fmt"
import "sync"

func TestHackerTarget(t *testing.T) {
	domain := "google.com"
	source := HackerTarget{}
	results := []string{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result.Success.(string))
	}

	// should be something like 4000 results
	if !(len(results) >= 1) {
		t.Errorf("expected to return more than one successful result, got %v", len(results))
	}
}

func TestHackerTarget_MultiThreaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := HackerTarget{}
	results := []string{}

	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			for result := range source.ProcessDomain(domain) {
				mx.Lock()
				results = append(results, result.Success.(string))
				mx.Unlock()
			}
		}(domain)
	}

	wg.Wait() // collect results

	if !(len(results) >= 4) {
		t.Errorf("expected over ( or exactly ) 4 results from multi-threaded example, got '%v'", len(results))
	}
}

func ExampleHackerTarget() {
	domain := "google.com"
	source := HackerTarget{}
	results := []string{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result.Success.(string))
	}

	fmt.Println(len(results) >= 1)
	// Output: true
}

func ExampleHackerTargetMultiThreaded() {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := HackerTarget{}
	results := []string{}

	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			for result := range source.ProcessDomain(domain) {
				mx.Lock()
				results = append(results, result.Success.(string))
				mx.Unlock()
			}
		}(domain)
	}

	wg.Wait() // collect results

	fmt.Println(len(results) >= 1)
	// Output: true
}
