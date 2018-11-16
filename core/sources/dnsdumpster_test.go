package sources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/subfinder/research/core"
)

func TestDNSDumpster(t *testing.T) {
	domain := "apple.com"
	source := DNSDumpster{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// stop after 20
	counter := 0

	for result := range source.ProcessDomain(ctx, domain) {
		counter++
		if counter == 20 {
			cancel()
		}
		t.Log(result.Success)
	}

	fmt.Println("found", counter, ctx.Err())
}

func TestDNSDumpster_Recursive(t *testing.T) {
	domain := "apple.com"
	source := &DNSDumpster{}
	results := []*core.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := &core.EnumerationOptions{
		Recursive: true,
		Sources:   []core.Source{source},
	}

	for result := range core.EnumerateSubdomains(ctx, domain, options) {
		results = append(results, result)
		fmt.Println(result)
	}

	fmt.Println(len(results), ctx.Err())
}
