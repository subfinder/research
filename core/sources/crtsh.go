package sources

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"

	"github.com/subfinder/research/core"
)

// CrtSh is a source to process subdomains from https://crt.sh
type CrtSh struct{}

type crtshObject struct {
	NameValue string `json:"name_value"`
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *CrtSh) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("crtsh", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		resp, err := core.HTTPClient.Get("https://crt.sh/?q=%25." + domain + "&output=json")
		if err != nil {
			results <- core.NewResult("crtsh", nil, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			results <- core.NewResult("crtsh", nil, errors.New(resp.Status))
			return
		}

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
					results <- core.NewResult("crtsh", nil, err)
					continue
				}
				// This could potentially be made more efficient.
				for _, str := range domainExtractor.FindAllString(object.NameValue, -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						results <- core.NewResult("crtsh", str, nil)
					}
				}
			}
		}
	}(domain, results)
	return results
}
