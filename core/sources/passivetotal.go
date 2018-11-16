package sources

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// Passivetotal is a source to process subdomains from https://passivetotal.org
type Passivetotal struct {
	APIToken    string
	APIUsername string
	lock        *semaphore.Weighted
}

type passivetotalObject struct {
	Subdomains []string `json:"subdomains"`
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Passivetotal) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	var resultLabel = "passivetotal"

	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		if source.APIToken == "" {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New("no api token")))
			return
		}

		if source.APIUsername == "" {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New("no api username")))
			return
		}

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		var body = []byte(`{"query":"` + domain + `"}`)

		req, err := http.NewRequest("GET", "https://api.passivetotal.org/v2/enrichment/subdomains", bytes.NewBuffer(body))

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		req.SetBasicAuth(source.APIUsername, source.APIToken)
		req.Header.Set("Content-Type", "application/json")

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

		hostResponse := passivetotalObject{}

		err = json.NewDecoder(resp.Body).Decode(&hostResponse)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		for _, sub := range hostResponse.Subdomains {
			str := sub + "." + domain
			if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
				return
			}
		}

	}(domain, results)
	return results
}
