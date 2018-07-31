package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"
import "strings"
import "errors"

// HackerTarget is a source to process subdomains from https://hackertarget.com
type HackerTarget struct {
	APIKey string
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *HackerTarget) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		httpClient := &http.Client{
			Timeout: time.Second * 4,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		}

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- &core.Result{Type: "hackertarget", Failure: err}
			return
		}

		uniqFilter := map[string]bool{}

		// get response from the API, optionally with an API key
		var resp *http.Response

		// check API key
		if source.APIKey != "" {
			resp, err = httpClient.Get("https://api.hackertarget.com/hostsearch/?q=" + domain + "&apikey=" + source.APIKey)
		} else {
			resp, err = httpClient.Get("https://api.hackertarget.com/hostsearch/?q=" + domain)
		}
		if err != nil {
			results <- &core.Result{Type: "hacker target", Failure: err}
			return
		}
		defer resp.Body.Close()

		// TODO: investigate io.LimitedReader
		// read response body, extracting subdomains
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			str := strings.Split(scanner.Text(), ",")[0]
			if strings.Contains(str, "API count exceeded") {
				results <- &core.Result{Type: "hackertarget", Failure: errors.New(str)}
			} else {
				for _, str := range domainExtractor.FindAllString(str, -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						results <- &core.Result{Type: "hackertarget", Success: str}
					}
				}
			}
		}
	}(domain, results)
	return results
}
