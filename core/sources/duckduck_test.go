package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestDuckDuckGo(t *testing.T) {
	domain := "google.com"
	source := DuckDuckGo{}
	results := []interface{}{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for result := range core.UniqResults(source.ProcessDomain(ctx, domain)) {
		results = append(results, result.Success)
	}

	fmt.Println(results)

	if !(len(results) >= 1) {
		t.Errorf("expected more than 1 result(s), got '%v'", len(results))
	}
}
