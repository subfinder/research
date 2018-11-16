package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestDNSTable(t *testing.T) {
	domain := "bing.com"
	source := DNSTable{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for result := range source.ProcessDomain(ctx, domain) {
		fmt.Println(result)
		results = append(results, result)
	}

	if !(len(results) >= 30) {
		t.Errorf("expected more than 30 result(s), got '%v'", len(results))
	}
}

func TestDNSTableRecursive(t *testing.T) {
	domain := "bing.com"
	source := &DNSTable{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	options := &core.EnumerationOptions{
		Recursive: true,
		Context:   ctx,
		Sources:   []core.Source{source},
	}

	for result := range core.EnumerateSubdomains(domain, options) {
		results = append(results, result)
		fmt.Println(result)
	}

	fmt.Println(len(results), ctx.Err())
}

// func TestFindSubdomainsDotComMultiThreaded(t *testing.T) {
// 	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
// 	source := FindSubdomainsDotCom{}
// 	results := []*core.Result{}
//
// 	wg := sync.WaitGroup{}
// 	mx := sync.Mutex{}
//
// 	for _, domain := range domains {
// 		wg.Add(1)
// 		go func(domain string) {
// 			defer wg.Done()
// 			for result := range source.ProcessDomain(domain) {
// 				mx.Lock()
// 				results = append(results, result)
// 				mx.Unlock()
// 			}
// 		}(domain)
// 	}
//
// 	wg.Wait() // collect results
//
// 	if len(results) <= 4000 {
// 		t.Errorf("expected at least 4000 results, got '%v'", len(results))
// 	}
// }
//
// func ExampleFindSubdomainsDotCom() {
// 	domain := "bing.com"
// 	source := FindSubdomainsDotCom{}
// 	results := []*core.Result{}
//
// 	for result := range source.ProcessDomain(domain) {
// 		results = append(results, result)
// 	}
//
// 	fmt.Println(len(results) >= 400)
// 	// Output: true
// }
//
// func ExampleFindSubdomainsDotCom_multiThreaded() {
// 	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
// 	source := FindSubdomainsDotCom{}
// 	results := []*core.Result{}
//
// 	wg := sync.WaitGroup{}
// 	mx := sync.Mutex{}
//
// 	for _, domain := range domains {
// 		wg.Add(1)
// 		go func(domain string) {
// 			defer wg.Done()
// 			for result := range source.ProcessDomain(domain) {
// 				mx.Lock()
// 				results = append(results, result)
// 				mx.Unlock()
// 			}
// 		}(domain)
// 	}
//
// 	wg.Wait() // collect results
//
// 	fmt.Println(len(results) >= 4000)
// 	// Output: true
// }
//
// func BenchmarkFindSubdomainsDotComSingleThreaded(b *testing.B) {
// 	domain := "bing.com"
// 	source := FindSubdomainsDotCom{}
//
// 	for n := 0; n < b.N; n++ {
// 		results := []*core.Result{}
// 		for result := range source.ProcessDomain(domain) {
// 			results = append(results, result)
// 		}
// 	}
// }
//
// func BenchmarkFindSubdomainsDotComMultiThreaded(b *testing.B) {
// 	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
// 	source := FindSubdomainsDotCom{}
// 	wg := sync.WaitGroup{}
// 	mx := sync.Mutex{}
//
// 	for n := 0; n < b.N; n++ {
// 		results := []*core.Result{}
//
// 		for _, domain := range domains {
// 			wg.Add(1)
// 			go func(domain string) {
// 				defer wg.Done()
// 				for result := range source.ProcessDomain(domain) {
// 					mx.Lock()
// 					results = append(results, result)
// 					mx.Unlock()
// 				}
// 			}(domain)
// 		}
//
// 		wg.Wait() // collect results
// 	}
// }
