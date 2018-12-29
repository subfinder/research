package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// Entrust is a source to process subdomains from https://entrust.com
type Entrust struct {
	lock *semaphore.Weighted
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Entrust) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(entrustLabel, nil, err))
			return
		}
		defer source.lock.Release(1)

		domainExtractor := core.NewSingleSubdomainExtractor(domain)

		req, err := http.NewRequest(http.MethodGet, "https://ctsearch.entrust.com/api/v1/certificates?fields=subjectDN&domain="+domain+"&includeExpired=true&exactMatch=false&limit=5000", nil)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(entrustLabel, nil, err))
			return
		}

		req.Cancel = ctx.Done()
		req.WithContext(ctx)

		resp, err := core.HTTPClient.Do(req)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(entrustLabel, nil, err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			sendResultWithContext(ctx, results, core.NewResult(entrustLabel, nil, errors.New(resp.Status)))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			str := domainExtractor([]byte(strings.Replace(scanner.Text(), "u003d", " ", -1)))
			if str != "" {
				if !sendResultWithContext(ctx, results, core.NewResult(entrustLabel, str, nil)) {
					return
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(entrustLabel, nil, err))
			return
		}

	}(domain, results)
	return results
}
