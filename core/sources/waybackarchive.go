package sources

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/subfinder/research/core"
)

// WaybackArchive is a source to process subdomains from http://web.archive.org
type WaybackArchive struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *WaybackArchive) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {

	var resultLabel = "waybackarchive"

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		uniqFilter := map[string]bool{}

		req, err := http.NewRequest(http.MethodGet, "http://web.archive.org/cdx/search/cdx?url=*."+domain+"/*&output=json&fl=original&collapse=urlkey", nil)
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
			if scanner.Bytes()[0] == 44 { // if ","
				str := string(jsonBuffer.Bytes())
				jsonBuffer.Reset()
				str = domainExtractor.FindString(str)
				_, found := uniqFilter[str]
				if !found && str != "" {
					uniqFilter[str] = true
					if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
						return
					}
				}
			} else {
				jsonBuffer.Write(scanner.Bytes())
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}
	}(domain, results)
	return results
}
