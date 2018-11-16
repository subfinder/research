package sources

import (
	"context"
	"runtime"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

func sendResultWithContext(ctx context.Context, results chan *core.Result, result *core.Result) bool {
	select {
	case <-ctx.Done():
		return false
	case results <- result:
		return true
	}
}

var maxWorkers = runtime.GOMAXPROCS(0)

func defaultLockValue() *semaphore.Weighted {
	return semaphore.NewWeighted(int64(maxWorkers))
}
