package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestArchiveIs(t *testing.T) {
	domain := "bing.com"
	source := ArchiveIs{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	//t.Log(len(results), "\n", "Success: ", success, "Failure: ", failure)
	if !(len(results) >= 20) {
		t.Errorf("expected more than 20 result(s), got '%v'", len(results))
	}
}

