package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/subfinder/research/core"
)

// Entrust is a source to process subdomains from https://entrust.com
type Entrust struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Entrust) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {

	var resultLabel = "entrust"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		uniqFilter := map[string]bool{}

		req, err := http.NewRequest(http.MethodGet, "https://ctsearch.entrust.com/api/v1/certificates?fields=subjectDN&domain="+domain+"&includeExpired=true&exactMatch=false&limit=5000", nil)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

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

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			txt := strings.Replace(scanner.Text(), "u003d", " ", -1)
			for _, str := range domainExtractor.FindAllString(txt, -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
						return
					}
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

	}(domain, results)
	return results
}
