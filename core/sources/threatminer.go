package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"

// Threatminer is a source to process subdomains from https://www.threatminer.org
type Threatminer struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Threatminer) ProcessDomain(domain string) <-chan *core.Result {
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
			results <- core.NewResult("threatminer", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := httpClient.Get("https://www.threatminer.org/getData.php?e=subdomains_container&q=" + domain + "&t=0&rt=10&p=1")
		if err != nil {
			results <- core.NewResult("threatminer", nil, err)
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- core.NewResult("threatminer", str, nil)
				}
			}
		}
	}(domain, results)
	return results
}
