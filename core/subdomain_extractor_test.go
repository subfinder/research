package core

import "testing"
import "fmt"

func TestNewSubdomainExtractor(t *testing.T) {
	_, err := NewSubdomainExtractor("google.com")
	if err != nil {
		t.Error(err)
	}
}

func ExampleNewSubdomainExtractor() {
	exampleText := `<a href="https://subdomain.google.com">`

	extractor, err := NewSubdomainExtractor("google.com")
	if err != nil {
		panic(err)
	}

	subdomain := extractor.FindString(exampleText)

	fmt.Println(subdomain)
	// Output: subdomain.google.com
}
