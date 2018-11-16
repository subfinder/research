package sources

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/subfinder/research/core"
)

// SecurityTrails is a source to process subdomains from https://securitytrails.com
type SecurityTrails struct {
	APIToken string
}

type securitytrailsObject struct {
	Subdomains []string `json:"subdomains"`
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *SecurityTrails) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {

	var resultLabel = "riddler"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

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
