package sources

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/subfinder/research/core"
)

// ArchiveIs is a source to process subdomains from http://archive.is
type ArchiveIs struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *ArchiveIs) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {

	var resultLabel = "archiveis"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		uniqFilter := map[string]bool{}

		for currentPage := 0; currentPage <= 750; currentPage += 10 {
			if ctx.Err() != nil {
				return
			}
			url := "https://archive.is/offset=" + strconv.Itoa(currentPage) + "/*." + domain

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

			scanner.Split(bufio.ScanWords)

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
