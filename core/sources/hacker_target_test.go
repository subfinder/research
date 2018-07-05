package sources

import "testing"
import "fmt"

func TestHackerTarget(t *testing.T) {
	domain := "google.com"
	source := HackerTarget{}
	results := []string{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result.Success.(string))
	}

	// should be something like 4000 results
	if !(len(results) >= 1) {
		t.Errorf("expected to return more than one successful result, got %v", len(results))
	}
}

