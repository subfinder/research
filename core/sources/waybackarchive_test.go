package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestWaybackArchive(t *testing.T) {
	domain := "bing.com"
	source := WaybackArchive{}
	results := []*core.Result{}

	success := 0
	failure := 0

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
		if result.Failure != nil {
			failure += 1
		} else {
			success += 1
		}
	}

	t.Log(len(results), "\n", "Success: ", success, "Failure: ", failure)
	if !(len(results) >= 1) {
		t.Errorf("expected more than 1 result(s), got '%v'", len(results))
	}
}

func TestWaybackArchiveMultiThreaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := WaybackArchive{}
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

	if len(results) <= 4 {
		t.Errorf("expected at least 4 results, got '%v'", len(results))
	}
}

