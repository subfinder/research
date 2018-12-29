package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// DNSTable is a source to process subdomains from https://dnstable.com
type DNSTable struct {
	lock *semaphore.Weighted
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *DNSTable) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(dnstableLabel, nil, err))
			return
		}
		defer source.lock.Release(1)

		domainExtractor := core.NewSingleSubdomainExtractor(domain)

		req, err := http.NewRequest(http.MethodGet, "https://dnstable.com/domain/"+domain, nil)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(dnstableLabel, nil, err))
			return
		}

		req.Cancel = ctx.Done()
		req.WithContext(ctx)

		resp, err := core.HTTPClient.Do(req)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(dnstableLabel, nil, err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			sendResultWithContext(ctx, results, core.NewResult(dnstableLabel, nil, errors.New(resp.Status)))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			str := domainExtractor(scanner.Bytes())
			if str != "" {
				if !sendResultWithContext(ctx, results, core.NewResult(dnstableLabel, str, nil)) {
					return
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(dnstableLabel, nil, err))
			return
		}

	}(domain, results)
	return results
}
