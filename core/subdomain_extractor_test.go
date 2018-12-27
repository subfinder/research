package core

import "testing"
import "fmt"
import "regexp"

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

func BenchmarkSubdomainExtractorOld(b *testing.B) {
	domain := "google.com"

	var result string

	extractor, err := regexp.Compile(`[a-zA-Z0-9\*_.-]+\.` + domain)
	if err != nil {
		b.Error(err)
	}

	for n := 0; n < b.N; n++ {
		exampleText := `<a href="https://subdomain.google.com">`

		result = extractor.FindString(exampleText)
	}

	b.Log(result)
}

func BenchmarkSubdomainExtractorNew(b *testing.B) {
	domain := "google.com"

	var result string

	extractor, err := regexp.Compile(`[\w-\*]+\.` + domain)
	if err != nil {
		b.Error(err)
	}

	for n := 0; n < b.N; n++ {

		exampleText := `<a href="https://subdomain.google.com">`

		result = extractor.FindString(exampleText)
	}

	b.Log(result)
}
