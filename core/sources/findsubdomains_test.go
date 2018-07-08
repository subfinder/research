package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestFindSubdomainsDotCom(t *testing.T) {
	domain := "bing.com"
	source := FindSubdomainsDotCom{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if !(len(results) >= 400) {
		t.Errorf("expected more than 400 result(s), got '%v'", len(results))
	}
}

