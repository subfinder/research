package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/subfinder/research/core"
)

// Bing is a source to process subdomains from https://bing.com
type Bing struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Bing) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	var resultLabel = "bing"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		uniqFilter := map[string]bool{}

		for currentPage := 1; currentPage <= 750; currentPage += 10 {
			url := "https://www.bing.com/search?q=domain%3A" + domain + "&go=Submit&first=" + strconv.Itoa(currentPage)
			req, err := http.NewRequest(http.MethodGet, url, nil)
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

			if resp.StatusCode != 200 {
				resp.Body.Close()
				sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New(resp.Status)))
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			for scanner.Scan() {
				for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
							resp.Body.Close()
							return
						}
					}
				}
			}

			resp.Body.Close()

		}

	}(domain, results)
	return results
}
