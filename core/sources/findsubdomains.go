package sources

import (
	"bufio"

	"github.com/subfinder/research/core"
)

// FindSubdomainsDotCom is a source to process subdomains from https://findsubdomains.com
type FindSubdomainsDotCom struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *FindSubdomainsDotCom) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("findsubdomainsdotcom", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := core.HTTPClient.Get("https://findsubdomains.com/subdomains-of/" + domain)
		if err != nil {
			results <- core.NewResult("findsubdomainsdotcom", nil, err)
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- core.NewResult("findsubdomainsdotcom", str, nil)
				}
			}
		}

	}(domain, results)
	return results
}
