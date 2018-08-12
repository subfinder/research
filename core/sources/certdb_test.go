package sources

import (
	"fmt"
	"sync"
	"testing"

	"github.com/subfinder/research/core"
)

func TestCertDB(t *testing.T) {
	domain := "google.com"
	source := CertDB{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		t.Log(result)
		results = append(results, result)
	}

	if !(len(results) >= 2) {
		t.Errorf("expected more than 20 result(s), got '%v'", len(results))
	}
}

func TestCertDB_multi_threaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := CertDB{}
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

func ExampleCertDB() {
	domain := "google.com"
	source := CertDB{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	fmt.Println(len(results) >= 2)
	// Output: true
}

func ExampleCertDB_multi_threaded() {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := CertDB{}
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

func BenchmarkCertDB_single_threaded(b *testing.B) {
	domain := "google.com"
	source := CertDB{}

	for n := 0; n < b.N; n++ {
		results := []*core.Result{}
		for result := range source.ProcessDomain(domain) {
			results = append(results, result)
		}
	}
}

func BenchmarkCertDB_multi_threaded(b *testing.B) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := CertDB{}
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
