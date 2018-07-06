package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestCrtSh(t *testing.T) {
	domain := "bing.com"
	source := CrtSh{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if !(len(results) >= 500) {
		t.Errorf("expected more than 500 result(s), got '%v'", len(results))
	}
}

