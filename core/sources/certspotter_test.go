package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestCertSpotter(t *testing.T) {
	domain := "google.com"
	source := CertSpotter{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uniqFilter := map[string]bool{}

	for result := range source.ProcessDomain(ctx, domain) {
		if result.IsSuccess() {
			str := result.Success.(string)
			_, found := uniqFilter[str]
			if !found {
				uniqFilter[str] = true
				fmt.Println(result)
				results = append(results, result)
			}
		}
	}

	if !(len(results) >= 3000) {
		t.Errorf("expected more than 3000 results, got '%v'", len(results))
	}
}

func TestCertSpotterRecursive(t *testing.T) {
	domain := "apple.com"
	source := &CertSpotter{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := &core.EnumerationOptions{
		Recursive: true,
		Context:   ctx,
		Cancel:    cancel,
		Sources:   []core.Source{source},
	}

	for result := range core.EnumerateSubdomains(domain, options) {
		results = append(results, result)
		fmt.Println(result)

	}

	if !(len(results) >= 5) {
		t.Errorf("expected more than 5 result(s), got '%v'", len(results))
		t.Error(ctx.Err())
	}
}

// func TestCertSpotter_MultiThreaded(t *testing.T) {
// 	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
// 	source := CertSpotter{}
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
// 				t.Log(result)
// 				mx.Lock()
// 				results = append(results, result)
// 				mx.Unlock()
// 			}
// 		}(domain)
// 	}
//
// 	wg.Wait() // collect results
//
// 	if len(results) < 6000 {
// 		t.Errorf("expected more than 6000 results, got '%v'", len(results))
// 	}
// }
//
// func ExampleCertSpotter() {
// 	domain := "google.com"
// 	source := CertSpotter{}
// 	results := []*core.Result{}
//
// 	for result := range source.ProcessDomain(domain) {
// 		results = append(results, result)
// 	}
//
// 	fmt.Println(len(results) >= 3000)
// 	// Output: true
// }
//
// func ExampleCertSpotter_multi_threaded() {
// 	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
// 	source := CertSpotter{}
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
// 	fmt.Println(len(results) > 6000)
// 	// Output: true
// }
//
// func BenchmarkCertSpotterSingleThreaded(b *testing.B) {
// 	domain := "google.com"
// 	source := CertSpotter{}
//
// 	for n := 0; n < b.N; n++ {
// 		results := []*core.Result{}
// 		for result := range source.ProcessDomain(domain) {
// 			results = append(results, result)
// 		}
// 	}
// }
//
// func BenchmarkCertSpotterMultiThreaded(b *testing.B) {
// 	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
// 	source := CertSpotter{}
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
