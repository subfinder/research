package sources

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/subfinder/research/core"
)

// CertSpotter is a source to process subdomains from https://certspotter.com
type CertSpotter struct {
	APIToken string
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *CertSpotter) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	var resultLabel = "certspotter"

	wg := sync.WaitGroup{}

	wg.Add(3)

	results := make(chan *core.Result)

	// apiv0
	go func(domain string, results chan *core.Result) {
		defer wg.Done()

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		url := "https://certspotter.com/api/v0/certs?domain=" + domain

		req, err := http.NewRequest(http.MethodGet, url, nil)

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

		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
					return
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

	}(domain, results)

	// apiv1 certs
	go func(domain string, results chan *core.Result) {
		defer wg.Done()

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		url := "https://api.certspotter.com/v1/certs?domain=" + domain + "&include_subdomains=true&expand=dns_names&match_wildcards=true"

		req, err := http.NewRequest(http.MethodGet, url, nil)

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		if source.APIToken != "" {
			req.Header.Set("Authorization", "Bearer "+source.APIToken)
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

		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}

			if !strings.Contains(scanner.Text(), `"data":"`) {
				continue
			}

			str := scanner.Text()
			str = strings.TrimSpace(str)
			str = strings.Trim(str, `"data":"`)
			str = strings.Trim(str, `"`)

			decodedData, err := base64.StdEncoding.DecodeString(str)

			if err != nil {
				continue
			}

			for _, str := range domainExtractor.FindAllString(string(decodedData), -1) {
				if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
					return
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

	}(domain, results)

	// apiv1 issuances
	go func(domain string, results chan *core.Result) {
		defer wg.Done()

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		url := "https://api.certspotter.com/v1/issuances?domain=" + domain + "&include_subdomains=true&expand=dns_names&match_wildcards=true"

		req, err := http.NewRequest(http.MethodGet, url, nil)

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

		if source.APIToken != "" {
			req.Header.Set("Authorization", "Bearer "+source.APIToken)
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

		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				if !sendResultWithContext(ctx, results, core.NewResult(resultLabel, str, nil)) {
					return
				}
			}
		}

		err = scanner.Err()

		if err != nil {
			sendResultWithContext(ctx, results, core.NewResult(resultLabel, nil, err))
			return
		}

	}(domain, results)

	go func() {
		defer close(results)

		wg.Wait()
	}()

	return results
}
