package sources

import (
	"bufio"
	"errors"

	"github.com/subfinder/research/core"
)

// CertDB is a source to process subdomains from https://certdb.com
type CertDB struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *ThreatCrowd) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("certdb", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := core.HTTPClient.Get("https://certdb.com/domain/" + domain)
		if err != nil {
			results <- core.NewResult("certdb", nil, err)
			return
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			results <- core.NewResult("certdb", nil, errors.New(resp.Status))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- core.NewResult("certdb", str, nil)
				}
			}
		}

		resp.Body.Close()

	}(domain, results)
	return results
}
