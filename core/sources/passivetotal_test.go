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
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for result := range source.ProcessDomain(ctx, domain) {
		fmt.Println(result)
		results = append(results, result)
	}

	fmt.Println("found", len(results), ctx.Err())
}
