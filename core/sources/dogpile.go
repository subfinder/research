package sources

import (
	"bufio"
	"strconv"

	"github.com/subfinder/research/core"
)

// DogPile is a source to process subdomains from http://dogpile.com
//
// Note
//
// This source uses http instead of https because of problems dogpile's SSL cert.
//
type DogPile struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *DogPile) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("dogpile", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		for currentPage := 1; currentPage <= 750; currentPage++ {
			resp, err := core.HTTPClient.Get("http://www.dogpile.com/search/web?q=" + domain + "&qsi=" + strconv.Itoa(currentPage*15+1))
			if err != nil {
				results <- core.NewResult("dogpile", nil, err)
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						results <- core.NewResult("dogpile", str, nil)
					}
				}
			}

			resp.Body.Close()

		}

	}(domain, results)
	return results
}
