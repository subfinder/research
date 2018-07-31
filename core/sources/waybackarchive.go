package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"
import "bytes"

// WaybackArchive is a source to process subdomains from http://web.archive.org
type WaybackArchive struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *WaybackArchive) ProcessDomain(domain string) <-chan *core.Result {
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
			results <- &core.Result{Type: "waybackarchive", Failure: err}
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := httpClient.Get("http://web.archive.org/cdx/search/cdx?url=*." + domain + "/*&output=json&fl=original&collapse=urlkey")
		if err != nil {
			results <- &core.Result{Type: "waybackarchive", Failure: err}
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		scanner.Split(bufio.ScanBytes)

		jsonBuffer := bytes.Buffer{}

		for scanner.Scan() {
			if scanner.Bytes()[0] == 44 { // if ","
				str := string(jsonBuffer.Bytes())
				jsonBuffer.Reset()
				str = domainExtractor.FindString(str)
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- &core.Result{Type: "waybackarchive", Success: str}
				}
			} else {
				jsonBuffer.Write(scanner.Bytes())
			}
		}
	}(domain, results)
	return results
}
