package sources

import (
	"bufio"
	"net"
	"net/http"
	"time"

	core "github.com/subfinder/research/core"
)

type DnsDbDotCom struct{}

func (source *DnsDbDotCom) ProcessDomain(domain string) <-chan *core.Result {
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
			results <- &core.Result{Type: "dnsdbdotcom", Failure: err}
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := httpClient.Get("http://www.dnsdb.org/f/" + domain + ".dnsdb.org/")
		if err != nil {
			results <- &core.Result{Type: "dnsdbdotcom", Failure: err}
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- &core.Result{Type: "dnsdbdotcom", Success: str}
				}
			}
		}

	}(domain, results)
	return results
}
