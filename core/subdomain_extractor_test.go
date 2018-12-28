package core

import "testing"
import "regexp"

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

	exampleText := `<a href="https://subdomain.google.com">`

	for n := 0; n < b.N; n++ {
		result = extractor.FindString(exampleText)
	}

	b.Log(result)
}

func BenchmarkSubdomainExtractorCustom(b *testing.B) {

	domain := "google.com"

	var result string

	exampleText := []byte(`<a href="https://subdomain.google.com">`)

	extractor := NewSingleSubdomainExtractor(domain)

	for n := 0; n < b.N; n++ {
		result = extractor(exampleText)
	}

	b.Log(result)
}

func BenchmarkSubdomainExtractorCustomMulti(b *testing.B) {

	domain := "google.com"

	var results []string

	exampleText := []byte(`<a href="https://subdomain.google.com">`)

	extractor := NewMultiSubdomainExtractor(domain)

	for n := 0; n < b.N; n++ {
		results = extractor(exampleText)
	}

	b.Log(results)
}
