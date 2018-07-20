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

	if !(len(results) >= 90) {
		t.Errorf("expected more than 90 result(s), got '%v'", len(results))
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

	if len(results) < 200 {
		t.Errorf("expected more than 200 results, got '%v'", len(results))
	}
}

func ExampleCrtSh() {
	domain := "bing.com"
	source := CrtSh{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	fmt.Println(len(results) >= 90)
	// Output: true
}

func ExampleCrtSh_multi_threaded() {
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

	fmt.Println(len(results) >= 200)
	// Output: true
}

func BenchmarkCrtShSingleThreaded(b *testing.B) {
	domain := "google.com"
	source := CrtSh{}

	for n := 0; n < b.N; n++ {
		results := []*core.Result{}
		for result := range source.ProcessDomain(domain) {
			results = append(results, result)
		}
	}
}

func BenchmarkCrtShMultiThreaded(b *testing.B) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := CrtSh{}
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
