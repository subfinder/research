package sources

import (
	"fmt"
	"sync"
	"testing"

	core "github.com/subfinder/research/core"
)

func TestDnsDbDotCom(t *testing.T) {
	domain := "bing.com"
	source := &DnsDbDotCom{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		t.Log(result)
		results = append(results, result)
	}

	if !(len(results) >= 400) {
		t.Errorf("expected more than 400 result(s), got '%v'", len(results))
	}
}

func TestDnsDbDotComMultiThreaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := &DnsDbDotCom{}
	results := []*core.Result{}

	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			for result := range source.ProcessDomain(domain) {
				t.Log(result)
				mx.Lock()
				results = append(results, result)
				mx.Unlock()
			}
		}(domain)
	}

	wg.Wait() // collect results

	if len(results) <= 350 {
		t.Errorf("expected at least 350 results, got '%v'", len(results))
	}
}

func ExampleDnsDbDotCom() {
	domain := "bing.com"
	source := &DnsDbDotCom{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	fmt.Println(len(results) >= 400)
	// Output: true
}

func ExampleDnsDbDotCom_multi_threaded() {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := &DnsDbDotCom{}
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

	fmt.Println(len(results) >= 350)
	// Output: true
}

func BenchmarkDnsDbDotComSingleThreaded(b *testing.B) {
	domain := "bing.com"
	source := &DnsDbDotCom{}

	for n := 0; n < b.N; n++ {
		results := []*core.Result{}
		for result := range source.ProcessDomain(domain) {
			results = append(results, result)
		}
	}
}

func BenchmarkDnsDbDotComMultiThreaded(b *testing.B) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := &DnsDbDotCom{}
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
