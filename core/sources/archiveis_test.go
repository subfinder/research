package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestArchiveIs(t *testing.T) {
	domain := "bing.com"
	source := ArchiveIs{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	//t.Log(len(results), "\n", "Success: ", success, "Failure: ", failure)
	if !(len(results) >= 20) {
		t.Errorf("expected more than 20 result(s), got '%v'", len(results))
	}
}

func TestArchiveIsMultiThreaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := ArchiveIs{}
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

	if len(results) <= 40 {
		t.Errorf("expected at least 40 results, got '%v'", len(results))
	}
}

func ExampleArchiveIs() {
	domain := "bing.com"
	source := ArchiveIs{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	fmt.Println(len(results) >= 20)
	// Output: true
}

