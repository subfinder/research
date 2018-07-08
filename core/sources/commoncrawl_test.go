package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestCommonCrawlDotOrg(t *testing.T) {
	domain := "bing.com"
	source := CommonCrawlDotOrg{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if !(len(results) >= 3) {
		t.Errorf("expected at least 3 result(s), got '%v'", len(results))
	}
}

