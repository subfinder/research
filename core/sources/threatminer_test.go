package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestThreatminer(t *testing.T) {
	domain := "bing.com"
	source := Threatminer{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if !(len(results) >= 140) {
		t.Errorf("expected more than 140 result(s), got '%v'", len(results))
	}
}

