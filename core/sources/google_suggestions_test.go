package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestGoogleSuggestions(t *testing.T) {
	domain := "google.com"
	source := GoogleSuggestions{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for result := range source.ProcessDomain(ctx, domain) {
		fmt.Println(result)
		results = append(results, result)
	}

	//if !(len(results) >= 400) {
	//	t.Errorf("expected more than 400 result(s), got '%v'", len(results))
	//}
}
