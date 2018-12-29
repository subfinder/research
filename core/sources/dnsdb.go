package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// DNSDbDotCom is a source to process subdomains from http://www.dnsdb.org/f/
type DNSDbDotCom struct {
	lock *semaphore.Weighted
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *DNSDbDotCom) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(dnsdbdLabel, nil, err))
			return
		}

		domainExtractor := core.NewSingleSubdomainExtractor(domain)

		defer source.lock.Release(1)

		req, err := http.NewRequest(http.MethodGet, "http://www.dnsdb.org/f/"+domain+".dnsdb.org/", nil)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(dnsdbdLabel, nil, err))
			return
		}

		req.Cancel = ctx.Done()
		req.WithContext(ctx)

		resp, err := core.HTTPClient.Do(req)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(dnsdbdLabel, nil, err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			sendResultWithContext(ctx, results, core.NewResult(dnsdbdLabel, nil, errors.New(resp.Status)))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			str := domainExtractor(scanner.Bytes())
			if str != "" {
				if !sendResultWithContext(ctx, results, core.NewResult(dnsdbdLabel, str, nil)) {
					return
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(dnsdbdLabel, nil, err))
			return
		}

	}(domain, results)
	return results
}
