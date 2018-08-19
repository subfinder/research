package sources

import (
	"bufio"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/subfinder/research/core"
)

// Baidu is a source to process subdomains from https://baidu.com
type Baidu struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Baidu) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("baidu", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		for currentPage := 1; currentPage <= 750; currentPage++ {
			resp, err := core.HTTPClient.Get("https://www.baidu.com/s?rn=10&pn=" + strconv.Itoa(currentPage) + "&wd=site%3A" + domain + "+-www.+&oq=site%3A" + domain + "+-www.+")
			if err != nil {
				results <- core.NewResult("baidu", nil, err)
				return
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				results <- core.NewResult("baidu", nil, errors.New(resp.Status))
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
						case <-ctx.Done():
							resp.Body.Close()
							return
						case results <- core.NewResult("baidu", str, nil):
							// move along
						}
					}
				}
			}

			resp.Body.Close()

		}

	}(domain, results)
	return results
}
