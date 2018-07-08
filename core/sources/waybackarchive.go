package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"
import "bytes"
import "regexp"
import "strings"

type WaybackArchive struct{}

func (source *WaybackArchive) ProcessDomain(domain string) <-chan *core.Result {
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

		resp, err := httpClient.Get("http://web.archive.org/cdx/search/cdx?url=*." + domain + "/*&output=json&fl=original&collapse=urlkey")
		if err != nil {
			results <- &core.Result{Type: "wayback archive", Failure: err}
			return
		}
		defer resp.Body.Close()

		domainExtractor, err := regexp.Compile(`(http|https)://\S+` + domain)
		if err != nil {
			results <- &core.Result{Type: "wayback archive", Failure: err}
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		scanner.Split(bufio.ScanBytes)

		jsonBuffer := bytes.Buffer{}

		uniqFilter := map[string]bool{}

		for scanner.Scan() {
			if scanner.Bytes()[0] == 44 { // if ","
				str := string(jsonBuffer.Bytes())
				jsonBuffer.Reset()
				str = domainExtractor.FindString(str)
				// a little extra finesse to message out the
				// actual subdomain from the string
				if str != "" {
					str = strings.Split(str, "://")[1]
					str = strings.Split(str, "/")[0]
					str = strings.Split(str, ":")[0]
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						results <- &core.Result{Type: "wayback archive", Success: str}
					}
				}
			} else {
				jsonBuffer.Write(scanner.Bytes())
			}
		}
	}(domain, results)
	return results
}
