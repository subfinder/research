package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"

// FindSubdomainsDotCom is a source to process subdomains from https://findsubdomains.com
type FindSubdomainsDotCom struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *FindSubdomainsDotCom) ProcessDomain(domain string) <-chan *core.Result {
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
			results <- &core.Result{Type: "findsubdomainsdotcom", Failure: err}
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := httpClient.Get("https://findsubdomains.com/subdomains-of/" + domain)
		if err != nil {
			results <- &core.Result{Type: "findsubdomainsdotcom", Failure: err}
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- &core.Result{Type: "findsubdomainsdotcom", Success: str}
				}
			}
		}

	}(domain, results)
	return results
}
