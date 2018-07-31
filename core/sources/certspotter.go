package sources

import (
	"bufio"
	"net"
	"net/http"
	"time"

	core "github.com/subfinder/research/core"
)

// CertSpotter is a source to process subdomains from https://certspotter.com
type CertSpotter struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *CertSpotter) ProcessDomain(domain string) <-chan *core.Result {
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
			results <- core.NewResult("certspotter", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		// get response from the API, optionally with an API key
		resp, err := httpClient.Get("https://certspotter.com/api/v0/certs?domain=" + domain)
		if err != nil {
			results <- core.NewResult("certspotter", nil, err)
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- core.NewResult("certspotter", str, nil)
				}
			}
		}

	}(domain, results)
	return results
}
