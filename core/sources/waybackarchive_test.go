package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestWaybackArchive(t *testing.T) {
	domain := "bing.com"
	source := WaybackArchive{}
	results := []*core.Result{}

	success := 0
	failure := 0

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
		if result.Failure != nil {
			failure += 1
		} else {
			success += 1
		}
	}

	t.Log(len(results), "\n", "Success: ", success, "Failure: ", failure)
	if !(len(results) >= 1) {
		t.Errorf("expected more than 1 result(s), got '%v'", len(results))
	}
}

