package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"
import "regexp"

type FindSubdomainsDotCom struct{}

func (source *FindSubdomainsDotCom) ProcessDomain(domain string) <-chan *core.Result {
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

		resp, err := httpClient.Get("https://findsubdomains.com/subdomains-of/" + domain)
		if err != nil {
			results <- &core.Result{Type: "findsubdomainsdotcom", Failure: err}
			return
		}
		defer resp.Body.Close()

		domainExtractor, err := regexp.Compile(`((\w|_|-|\*)+\.)+` + domain)
		if err != nil {
			results <- &core.Result{Type: "findsubdomainsdotcom", Failure: err}
			return
		}

		uniqFilter := map[string]bool{}

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
