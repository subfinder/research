package sources

import (
	"bufio"
	"errors"
	"strconv"

	"github.com/subfinder/research/core"
)

// Bing is a source to process subdomains from https://bing.com
type Bing struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Bing) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("bing", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		for currentPage := 1; currentPage <= 750; currentPage += 10 {
			resp, err := core.HTTPClient.Get("https://www.bing.com/search?q=domain%3A" + domain + "&go=Submit&first=" + strconv.Itoa(currentPage))
			if err != nil {
				results <- core.NewResult("bing", nil, err)
				return
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				results <- core.NewResult("bing", nil, errors.New(resp.Status))
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			for scanner.Scan() {
				for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						results <- core.NewResult("bing", str, nil)
					}
				}
			}

			resp.Body.Close()

		}

	}(domain, results)
	return results
}
