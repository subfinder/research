package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "encoding/json"
import "bufio"
import "bytes"

type CrtSh struct{}

type crtshObject struct {
	NameValue string `json:"name_value"`
}

func (source *CrtSh) ProcessDomain(domain string) <-chan *core.Result {
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

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- &core.Result{Type: "crtsh", Failure: err}
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := httpClient.Get("https://crt.sh/?q=%25." + domain + "&output=json")
		if err != nil {
			results <- &core.Result{Type: "crtsh", Failure: err}
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		scanner.Split(bufio.ScanBytes)

		jsonBuffer := bytes.Buffer{}
		for scanner.Scan() {
			jsonBuffer.Write(scanner.Bytes())
			if scanner.Bytes()[0] == 125 { // if "}"
				object := &crtshObject{}
				json.Unmarshal(jsonBuffer.Bytes(), &object)
				err = json.Unmarshal(jsonBuffer.Bytes(), &object)
				jsonBuffer.Reset()
				if err != nil {
					results <- &core.Result{Type: "crtsh", Failure: err}
					continue
				}
				// This could potentially be made more efficient.
				for _, str := range domainExtractor.FindAllString(object.NameValue, -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						results <- &core.Result{Type: "certspotter", Success: str}
					}
				}
			}
		}
	}(domain, results)
	return results
}
