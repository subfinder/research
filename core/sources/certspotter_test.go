package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestCertSpotter(t *testing.T) {
	domain := "google.com"
	source := CertSpotter{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		if result.IsFailure() {
			t.Fatal(result.Failure)
		}
		results = append(results, result)
	}

	if !(len(results) >= 5000) {
		t.Errorf("expected more than 5000 results, got '%v'", len(results))
	}
}

