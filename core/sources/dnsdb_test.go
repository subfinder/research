package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestDNSDbDotCom(t *testing.T) {
	domain := "bing.com"
	source := &DNSDbDotCom{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for result := range source.ProcessDomain(ctx, domain) {
		fmt.Println(result)
		results = append(results, result)
	}

	if !(len(results) >= 400) {
		t.Errorf("expected more than 400 result(s), got '%v'", len(results))
	}
}

func TestDNSDbDotComRecursive(t *testing.T) {
	domain := "bing.com"
	source := &DNSDbDotCom{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	options := &core.EnumerationOptions{
		Recursive: true,
		Sources:   []core.Source{source},
	}

	for result := range core.EnumerateSubdomains(ctx, domain, options) {
		results = append(results, result)
		fmt.Println(result)

	}

	if !(len(results) >= 5) {
		t.Errorf("expected more than 5 result(s), got '%v'", len(results))
		t.Error(ctx.Err())
	}

	fmt.Println(len(results), ctx.Err())
}

//func TestDNSDbDotComMultiThreaded(t *testing.T) {
//	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
//	source := &DNSDbDotCom{}
//	results := []*core.Result{}
//
//	wg := sync.WaitGroup{}
//	mx := sync.Mutex{}
//
//	for _, domain := range domains {
//		wg.Add(1)
//		go func(domain string) {
//			defer wg.Done()
//			for result := range source.ProcessDomain(domain) {
//				t.Log(result)
//				mx.Lock()
//				results = append(results, result)
//				mx.Unlock()
//			}
//		}(domain)
//	}
//
//	wg.Wait() // collect results
//
//	if len(results) <= 350 {
//		t.Errorf("expected at least 350 results, got '%v'", len(results))
//	}
//}
//
//func ExampleDNSDbDotCom() {
//	domain := "bing.com"
//	source := &DNSDbDotCom{}
//	results := []*core.Result{}
//
//	for result := range source.ProcessDomain(domain) {
//		results = append(results, result)
//	}
//
//	fmt.Println(len(results) >= 400)
//	// Output: true
//}
//
//func ExampleDNSDbDotCom_multi_threaded() {
//	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
//	source := &DNSDbDotCom{}
//	results := []*core.Result{}
//
//	wg := sync.WaitGroup{}
//	mx := sync.Mutex{}
//
//	for _, domain := range domains {
//		wg.Add(1)
//		go func(domain string) {
//			defer wg.Done()
//			for result := range source.ProcessDomain(domain) {
//				mx.Lock()
//				results = append(results, result)
//				mx.Unlock()
//			}
//		}(domain)
//	}
//
//	wg.Wait() // collect results
//
//	fmt.Println(len(results) >= 350)
//	// Output: true
//}
//
//func BenchmarkDNSDbDotComSingleThreaded(b *testing.B) {
//	domain := "bing.com"
//	source := &DNSDbDotCom{}
//
//	for n := 0; n < b.N; n++ {
//		results := []*core.Result{}
//		for result := range source.ProcessDomain(domain) {
//			results = append(results, result)
//		}
//	}
//}
//
//func BenchmarkDNSDbDotComMultiThreaded(b *testing.B) {
//	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
//	source := &DNSDbDotCom{}
//	wg := sync.WaitGroup{}
//	mx := sync.Mutex{}
//
//	for n := 0; n < b.N; n++ {
//		results := []*core.Result{}
//
//		for _, domain := range domains {
//			wg.Add(1)
//			go func(domain string) {
//				defer wg.Done()
//				for result := range source.ProcessDomain(domain) {
//					mx.Lock()
//					results = append(results, result)
//					mx.Unlock()
//				}
//			}(domain)
//		}
//
//		wg.Wait() // collect results
//	}
//}
