package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"
import "strings"
import "errors"

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
	resp, err := httpClient.Get("https://api.hackertarget.com/hostsearch/?q=")
	if err != nil {
		return true
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "API count exceeded") {
			return true
		}
		break
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
			str := strings.Split(scanner.Text(), ",")[0]
			if strings.Contains(str, "API count exceeded") {
				results <- &core.Result{Type: "hacker target", Failure: errors.New(str)}
			} else {
				results <- &core.Result{Type: "hacker target", Success: str}
			}
		}
	}(domain, results)
	return results
}
