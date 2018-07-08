package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestFindSubdomainsDotCom(t *testing.T) {
	domain := "bing.com"
	source := FindSubdomainsDotCom{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if !(len(results) >= 400) {
		t.Errorf("expected more than 400 result(s), got '%v'", len(results))
	}
}

func TestFindSubdomainsDotComMultiThreaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := FindSubdomainsDotCom{}
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

	if len(results) <= 4000 {
		t.Errorf("expected at least 4000 results, got '%v'", len(results))
	}
}

