package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestAsk(t *testing.T) {
	domain := "google.com"
	source := Ask{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for result := range source.ProcessDomain(ctx, domain) {
		//fmt.Println(result)
		results = append(results, result)
		fmt.Println(result)
		// Not waiting around to iterate all the possible pages.
		//if len(results) >= 5 {
		//	cancel()
		//}
	}

	if !(len(results) >= 5) {
		t.Errorf("expected more than 5 result(s), got '%v'", len(results))
		t.Error(ctx.Err())
	}
}

func TestAskRecursive(t *testing.T) {
	domain := "google.com"
	source := &Ask{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

//func TestAsk_multi_threaded(t *testing.T) {
//	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
//	source := Ask{}
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
//				fmt.Println(result)
//				mx.Lock()
//				results = append(results, result)
//				mx.Unlock()
//			}
//		}(domain)
//	}
//
//	wg.Wait() // collect results
//
//	if len(results) <= 4 {
//		t.Errorf("expected at least 4 results, got '%v'", len(results))
//	}
//}
//
//func ExampleAsk() {
//	domain := "google.com"
//	source := Ask{}
//	results := []*core.Result{}
//
//	for result := range source.ProcessDomain(domain) {
//		results = append(results, result)
//	}
//
//	fmt.Println(len(results) >= 20)
//	// Output: true
//}
//
//func ExampleAsk_multi_threaded() {
//	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
//	source := Ask{}
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
//	fmt.Println(len(results) >= 4)
//	// Output: true
//}
//
//func BenchmarkAsk_single_threaded(b *testing.B) {
//	domain := "google.com"
//	source := Ask{}
//
//	for n := 0; n < b.N; n++ {
//		results := []*core.Result{}
//		for result := range source.ProcessDomain(domain) {
//			results = append(results, result)
//		}
//	}
//}
//
//func BenchmarkAsk_multi_threaded(b *testing.B) {
//	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
//	source := Ask{}
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
