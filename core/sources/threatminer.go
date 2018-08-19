package sources

import (
	"bufio"
	"context"
	"errors"
	"time"

	"github.com/subfinder/research/core"
)

// Threatminer is a source to process subdomains from https://www.threatminer.org
type Threatminer struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Threatminer) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("threatminer", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := core.HTTPClient.Get("https://www.threatminer.org/getData.php?e=subdomains_container&q=" + domain + "&t=0&rt=10&p=1")
		if err != nil {
			results <- core.NewResult("threatminer", nil, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			results <- core.NewResult("archiveis", nil, errors.New(resp.Status))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		scanner.Split(bufio.ScanWords)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					select {
					case results <- core.NewResult("certdb", str, nil):
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
