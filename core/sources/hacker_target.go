package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"
import "strings"

type HackerTarget struct{}

func (source *HackerTarget) IsOverFreeLimit() bool {
	httpClient := &http.Client{
		Timeout: time.Second * 4,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	// get response from the API
	resp, err := httpClient.Get("https://api.hackertarget.com/hostsearch/?q=" + domain)
	if err != nil {
		results <- &core.Result{Type: "hacker target", Failure: err}
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "API count exceeded") {
			return true
		}
	}
	return false
}

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

		// get response from the API
		resp, err := httpClient.Get("https://api.hackertarget.com/hostsearch/?q=" + domain)
		if err != nil {
			results <- &core.Result{Type: "hacker target", Failure: err}
			return
		}
		defer resp.Body.Close()

		// TODO: investigate io.LimitedReader
		// read response body, extracting subdomains
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			results <- &core.Result{Type: "hacker target", Success: strings.Split(scanner.Text(), ",")[0]}
		}
	}(domain, results)
	return results
}
