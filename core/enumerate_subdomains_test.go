package core

import "testing"
import "fmt"

type FakeSource1 struct{}
type FakeSource2 struct{}

func (s *FakeSource1) ProcessDomain(domain string) <-chan *Result {
	results := make(chan *Result)
	go func(domain string) {
		defer close(results)
		for _, subdomain := range []string{"www.", "info.", "wiggle."} {
			results <- &Result{Success: subdomain + domain}
		}
	}(domain)
	return results
}

func (s *FakeSource2) ProcessDomain(domain string) <-chan *Result {
	results := make(chan *Result)
	go func(domain string) {
		defer close(results)
		for _, subdomain := range []string{"admin.", "user.", "mod."} {
			results <- &Result{Success: subdomain + domain}
		}
	}(domain)
	return results
}

func TestEnumerateSubdomains(t *testing.T) {
	domain := "google.com"
	options := &EnumerationOptions{Sources: []Source{&FakeSource1{}, &FakeSource2{}}}

	collectedResults := []string{}
	for result := range EnumerateSubdomains(domain, options) {
		collectedResults = append(collectedResults, result.Success.(string))
	}

	if len(collectedResults) != 6 {
		t.Error("unable to successfully enumerate through all subdomain sources")
	}
}

func ExampleEnumerateSubdomains() {
	domain := "google.com"

	sources := []Source{
		&FakeSource1{},
		&FakeSource2{},
	}

	options := &EnumerationOptions{Sources: sources}

	counter := 0

	for result := range EnumerateSubdomains(domain, options) {
		if result.IsSuccess() {
			counter++
		}
	}

	fmt.Println(counter)
	// Output: 6
}
