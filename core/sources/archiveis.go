package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"
import "regexp"
import "strings"

type ArchiveIs struct{}

func (source *ArchiveIs) ProcessDomain(domain string) <-chan *core.Result {
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

		resp, err := httpClient.Get("http://archive.is/*." + domain)
		if err != nil {
			results <- &core.Result{Type: "archiveis", Failure: err}
			return
		}
		defer resp.Body.Close()

		domainExtractor, err := regexp.Compile(`(http|https)://((\w|_|-|\*)+\.)+` + domain)
		if err != nil {
			results <- &core.Result{Type: "archiveis", Failure: err}
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				str = strings.Split(str, "://")[1]
				results <- &core.Result{Type: "archiveis", Success: str}
			}
		}
	}(domain, results)
	return results
}
