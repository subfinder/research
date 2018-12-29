package sources

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// GoogleSuggestions is a source to process subdomains from https://suggestqueries.google.com
type GoogleSuggestions struct {
	lock *semaphore.Weighted
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *GoogleSuggestions) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	var resultLabel = "google-suggestions"

	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}
		defer source.lock.Release(1)

		domainExtractor := core.NewSingleSubdomainExtractor(domain)

		req, err := http.NewRequest(http.MethodGet, "https://www.google.com/complete/search?output=search&client=chrome&q="+domain, nil)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		req.Cancel = ctx.Done()
		req.WithContext(ctx)

		resp, err := core.HTTPClient.Do(req)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New(resp.Status)))
			return
		}

		raw := []json.RawMessage{}

		err = json.NewDecoder(resp.Body).Decode(&raw)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		if len(raw) < 2 {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New("no suggestion data found")))
			return
		}

		sgs := []string{}

		err = json.Unmarshal(raw[1], &sgs)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		for _, s := range sgs {
			if ctx.Err() != nil {
				return
			}
			str := domainExtractor([]byte(s))
			if str != "" {
				if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
					return
				}
			}
		}
	}(domain, results)
	return results
}
