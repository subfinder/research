package sources

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/subfinder/research/core"
)

// CrtSh is a source to process subdomains from https://crt.sh
type CrtSh struct{}

type crtshObject struct {
	NameValue string `json:"name_value"`
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *CrtSh) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {

	var resultLabel = "crtsh"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		uniqFilter := map[string]bool{}

		req, err := http.NewRequest(http.MethodGet, "https://crt.sh/?q=%25."+domain+"&output=json", nil)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		req.WithContext(ctx)

		resp, err := core.HTTPClient.Do(req)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, errors.New(resp.Status)))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		scanner.Split(bufio.ScanBytes)

		jsonBuffer := bytes.Buffer{}
		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			jsonBuffer.Write(scanner.Bytes())
			if scanner.Bytes()[0] == 125 { // if "}"
				object := &crtshObject{}
				json.Unmarshal(jsonBuffer.Bytes(), &object)
				err = json.Unmarshal(jsonBuffer.Bytes(), &object)
				jsonBuffer.Reset()
				if err != nil {
					sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
					continue
				}
				// This could potentially be made more efficient.
				for _, str := range domainExtractor.FindAllString(object.NameValue, -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
							return
						}
					}
				}
			}

			err = scanner.Err()

			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
				return
			}
		}
	}(domain, results)
	return results
}
