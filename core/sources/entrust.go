package sources

import (
	"bufio"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/subfinder/research/core"
)

// Entrust is a source to process subdomains from https://entrust.com
type Entrust struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Entrust) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("entrust", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := core.HTTPClient.Get("https://ctsearch.entrust.com/api/v1/certificates?fields=subjectDN&domain=" + domain + "&includeExpired=true&exactMatch=false&limit=5000")
		if err != nil {
			results <- core.NewResult("entrust", nil, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			resp.Body.Close()
			results <- core.NewResult("entrust", nil, errors.New(resp.Status))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for scanner.Scan() {
			txt := strings.Replace(scanner.Text(), "u003d", " ", -1)
			for _, str := range domainExtractor.FindAllString(txt, -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					select {
					case results <- core.NewResult("entrust", str, nil):
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
