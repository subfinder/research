package sources

import (
	"context"

	"github.com/subfinder/research/core"
)

func sendResultWithContext(ctx context.Context, results chan *core.Result, result *core.Result) bool {
	select {
	case <-ctx.Done():
		return false
	case results <- result:
		return true
	}
}
