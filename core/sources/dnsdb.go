package sources

import (
	"bufio"
	"net"
	"net/http"
	"regexp"
	"strings"
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

		resp, err := httpClient.Get("http://www.dnsdb.org/f/" + domain + ".dnsdb.org/")
		if err != nil {
			results <- &core.Result{Type: "dnsdbdotcom", Failure: err}
			return
		}
		defer resp.Body.Close()

		reDomains := regexp.MustCompile("<a[^>]*?[^>]*>(.*?)</a>")

		uniqFilter := map[string]bool{}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range reDomains.FindAllString(scanner.Text(), -1) {
				str = strings.TrimRight(strings.Split(str, "\">")[1], "</a>")
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
