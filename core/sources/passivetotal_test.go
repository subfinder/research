package sources

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestPassivetotal(t *testing.T) {
	domain := "bing.com"
	source := Passivetotal{
		APIToken:    os.Getenv("PassivetotalKey"),
		APIUsername: os.Getenv("PassivetotalUsername"),
	}
	results := []interface{}{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for result := range core.UniqResults(source.ProcessDomain(ctx, domain)) {
		results = append(results, result.Success)
	}

	fmt.Println(results)

	fmt.Println("found", len(results), ctx.Err())
}
