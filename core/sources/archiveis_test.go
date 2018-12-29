package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestArchiveIs(t *testing.T) {
	domain := "apple.com"
	source := ArchiveIs{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// stop after 20
	counter := 0

	for result := range core.UniqResults(source.ProcessDomain(ctx, domain)) {
		counter++
		if counter == 20 {
			cancel()
		}
		fmt.Println(result.Success)
	}

	fmt.Println("found", counter, ctx.Err())
}

func TestArchiveIsRecursive(t *testing.T) {
	domain := "apple.com"
	source := &ArchiveIs{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := &core.EnumerationOptions{
		Recursive: true,
		Sources:   []core.Source{source},
	}

	for result := range core.UniqResults(core.EnumerateSubdomains(ctx, domain, options)) {
		results = append(results, result)
		fmt.Println(result)
	}

	fmt.Println(len(results), ctx.Err())
}
