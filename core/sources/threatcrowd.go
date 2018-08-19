package sources

import (
	"bufio"
	"context"
	"errors"
	"time"

	"github.com/subfinder/research/core"
)

// ThreatCrowd is a source to process subdomains from https://threatcrowd.com
type ThreatCrowd struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *ThreatCrowd) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("threatcrowd", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := core.HTTPClient.Get("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=" + domain)
		if err != nil {
			results <- core.NewResult("threatcrowd", nil, err)
			return
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			results <- core.NewResult("threatcrowd", nil, errors.New(resp.Status))
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
					case results <- core.NewResult("threatcrowd", str, nil):
						// move along
					case <-ctx.Done():
						resp.Body.Close()
						return
					}
				}
			}
		}

		resp.Body.Close()

	}(domain, results)
	return results
}
