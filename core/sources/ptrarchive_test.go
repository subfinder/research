package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestPTRArchiveDotCom(t *testing.T) {
	domain := "bing.com"
	source := PTRArchiveDotCom{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if !(len(results) >= 35) {
		t.Errorf("expected more than 35 result(s), got '%v'", len(results))
	}
}

func TestPTRArchiveDotComMultiThreaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := PTRArchiveDotCom{}
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

	if len(results) <= 500 {
		t.Errorf("expected at least 500 results, got '%v'", len(results))
	}
}

func ExamplePTRArchiveDotCom() {
	domain := "bing.com"
	source := PTRArchiveDotCom{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	fmt.Println(len(results) >= 35)
	// Output: true
}

func ExamplePTRArchiveDotComMultiThreaded() {
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

	fmt.Println(len(results) >= 500)
	// Output: true
}

func BenchmarkPTRArchiveDotComSingleThreaded(b *testing.B) {
	domain := "bing.com"
	source := FindSubdomainsDotCom{}

	for n := 0; n < b.N; n++ {
		results := []*core.Result{}
		for result := range source.ProcessDomain(domain) {
			results = append(results, result)
		}
	}
}

func BenchmarkPTRArchiveDotComMultiThreaded(b *testing.B) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := FindSubdomainsDotCom{}
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
