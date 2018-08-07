package sources

import (
	"bufio"

	"github.com/subfinder/research/core"
)

// DnsDbDotCom is a source to process subdomains from http://www.dnsdb.org/f/
type DnsDbDotCom struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *DnsDbDotCom) ProcessDomain(domain string) <-chan *core.Result {
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

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("dnsdbdotcom", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := httpClient.Get("http://www.dnsdb.org/f/" + domain + ".dnsdb.org/")
		if err != nil {
			results <- core.NewResult("dnsdbdotcom", nil, err)
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- core.NewResult("dnsdbdotcom", str, nil)
				}
			}
		}

	}(domain, results)
	return results
}
