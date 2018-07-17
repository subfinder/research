package sources

import (
	"bufio"
	"net"
	"net/http"
	"regexp"
	"time"

	core "github.com/subfinder/research/core"
)

type PTRArchiveDotCom struct{}

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

		resp, err := httpClient.Get("http://ptrarchive.com/tools/search3.htm?label=" + domain + "&date=ALL")
		if err != nil {
			results <- &core.Result{Type: "ptrarchivedotcom", Failure: err}
			return
		}
		defer resp.Body.Close()

		domainExtractor, err := regexp.Compile("] (.*) \\[")
		if err != nil {
			results <- &core.Result{Type: "ptrarchivedotcom", Failure: err}
			return
		}

		uniqFilter := map[string]bool{}

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
