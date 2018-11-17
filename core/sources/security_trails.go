package sources

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// SecurityTrails is a source to process subdomains from https://securitytrails.com
type SecurityTrails struct {
	APIToken string
	lock     *semaphore.Weighted
}

type securitytrailsObject struct {
	Subdomains []string `json:"subdomains"`
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *SecurityTrails) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	var resultLabel = "riddler"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}
		defer source.lock.Release(1)

		// check if only password was given
		if source.APIToken == "" {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New("no api token")))
			return
		}

		url := "https://api.securitytrails.com/v1/domain/" + domain + "/subdomains"
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		req.Header.Add("APIKEY", source.APIToken)

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

		hostResponse := securitytrailsObject{}

		err = json.NewDecoder(resp.Body).Decode(&hostResponse)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		for _, sub := range hostResponse.Subdomains {
			str := sub + "." + domain
			if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
				break
			}
		}
		return

	}(domain, results)
	return results
}
