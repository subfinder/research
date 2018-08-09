package sources

import (
	"bufio"
	"errors"
	"strconv"

	"github.com/subfinder/research/core"
)

// Yahoo is a source to process subdomains from https://yahoo.com
type Yahoo struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Yahoo) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("yahoo", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		for currentPage := 1; currentPage <= 750; currentPage++ {
			resp, err := core.HTTPClient.Get("https://search.yahoo.com/search?p=site:" + domain + "&b=" + strconv.Itoa(currentPage*10) + "&pz=10&bct=0&xargs=0")
			if err != nil {
				results <- core.NewResult("yahoo", nil, err)
				return
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				results <- core.NewResult("yahoo", nil, errors.New(resp.Status))
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						results <- core.NewResult("yahoo", str, nil)
					}
				}
			}

			resp.Body.Close()

		}

	}(domain, results)
	return results
}
