package sources

import (
	"bufio"
	"context"
	"errors"
	"time"

	"github.com/subfinder/research/core"
)

// CommonCrawlDotOrg is a source to process subdomains from http://commoncrawl.org
type CommonCrawlDotOrg struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *CommonCrawlDotOrg) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("commoncrawldotorg", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := core.HTTPClient.Get("https://index.commoncrawl.org/CC-MAIN-2018-17-index?url=*." + domain + "&output=json")
		if err != nil {
			results <- core.NewResult("commoncrawldotorg", nil, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			results <- core.NewResult("commoncrawldotorg", nil, errors.New(resp.Status))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					select {
					case results <- core.NewResult("commoncrawldotorg", str, nil):
					case results <- core.NewResult("bing", str, nil):
						// move along
					case <-ctx.Done():
						resp.Body.Close()
						return
					}
				}
			}
		}

	}(domain, results)
	return results
}
