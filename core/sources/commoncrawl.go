package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"

// CommonCrawlDotOrg is a source to process subdomains from http://commoncrawl.org
type CommonCrawlDotOrg struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *CommonCrawlDotOrg) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		httpClient := &http.Client{
			//Timeout: time.Second * 60,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		}

		resp, err := httpClient.Get("http://index.commoncrawl.org/CC-MAIN-2018-17-index?url=*." + domain + "&output=json")
		if err != nil {
			results <- core.NewResult("commoncrawldotorg", nil, err)
			return
		}
		defer resp.Body.Close()

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("commoncrawldotorg", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- core.NewResult("commoncrawldotorg", str, nil)
				}
			}
		}

	}(domain, results)
	return results
}
