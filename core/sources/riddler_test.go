package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestRiddler(t *testing.T) {
	domain := "bing.com"
	source := Riddler{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if !(len(results) >= 9) {
		t.Errorf("expected more than 9 result(s), got '%v'", len(results))
	}
}

