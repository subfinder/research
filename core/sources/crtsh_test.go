package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestCrtSh(t *testing.T) {
	domain := "bing.com"
	source := CrtSh{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if !(len(results) >= 500) {
		t.Errorf("expected more than 500 result(s), got '%v'", len(results))
	}
}

func TestCrtSh_MultiThreaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := CrtSh{}
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

	if len(results) < 1180 {
		t.Errorf("expected more than 1180 results, got '%v'", len(results))
	}
}

func ExampleCrtSh() {
	domain := "bing.com"
	source := CrtSh{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	fmt.Println(len(results) >= 500)
	// Output: true
}

func ExampleCrtSh_MultiThreaded() {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := CrtSh{}
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

	fmt.Println(len(results) > 1180)
	// Output: true
}

