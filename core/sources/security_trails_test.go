package sources

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSecurityTrails(t *testing.T) {
	domain := "apple.com"
	source := SecurityTrails{APIToken: os.Getenv("SecurityTrailsKey")}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// stop after 20
	counter := 0

	for result := range source.ProcessDomain(ctx, domain) {
		counter++
		if counter == 20 {
			cancel()
		}
		fmt.Println(result.Success)
	}

	fmt.Println("found", counter, ctx.Err())
}
