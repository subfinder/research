package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"
import "regexp"

type CertDB struct{}

func (source *CertDB) ProcessDomain(domain string) <-chan *core.Result {
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

		uniqFilter := map[string]bool{}

		domainExtractor, err := regexp.Compile(`[a-zA-Z0-9\*_.-]+\.` + domain)
		if err != nil {
			results <- &core.Result{Type: "certdb", Failure: err}
			return
		}

		resp, err := httpClient.Get("https://certdb.com/domain/" + domain)
		if err != nil {
			results <- &core.Result{Type: "certdb", Failure: err}
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- &core.Result{Type: "certdb", Success: str}
				}
			}
		}
	}(domain, results)
	return results
}
