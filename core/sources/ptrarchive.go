package sources

import (
	"bufio"
	"net"
	"net/http"
	"time"

	core "github.com/subfinder/research/core"
)

// PTRArchiveDotCom is a source to process subdomains from http://ptrarchive.com/
type PTRArchiveDotCom struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *PTRArchiveDotCom) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		httpClient := &http.Client{
			Timeout: time.Second * 60,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		}

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- &core.Result{Type: "ptrarchivedotcom", Failure: err}
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := httpClient.Get("http://ptrarchive.com/tools/search3.htm?label=" + domain + "&date=ALL")
		if err != nil {
			results <- &core.Result{Type: "ptrarchivedotcom", Failure: err}
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- &core.Result{Type: "ptrarchivedotcom", Success: str}
				}
			}
		}
	}(domain, results)
	return results
}
