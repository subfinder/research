package sources

import (
	"bufio"
	"context"
	"errors"

	"golang.org/x/sync/semaphore"
	"net/http"
	"strconv"

	"github.com/subfinder/research/core"
)

// ArchiveIs is a source to process subdomains from http://archive.is
type ArchiveIs struct {
	lock *semaphore.Weighted
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *ArchiveIs) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(archiveisLabel, nil, err))
			return
		}
		defer source.lock.Release(1)

		domainExtractor := core.NewSingleSubdomainExtractor(domain)

		for currentPage := 0; currentPage <= 750; currentPage += 10 {
			if ctx.Err() != nil {
				return
			}

			url := "https://archive.is/offset=" + strconv.Itoa(currentPage) + "/*." + domain

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(archiveisLabel, nil, err))
				return
			}

			req.Cancel = ctx.Done()
			req.WithContext(ctx)

			resp, err := core.HTTPClient.Do(req)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(archiveisLabel, nil, err))
				return
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				sendResultWithContext(ctx, results, core.NewResult(archiveisLabel, nil, errors.New(resp.Status)))
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				if ctx.Err() != nil {
					return
				}

				str := domainExtractor(scanner.Bytes())

				if !sendResultWithContext(ctx, results, core.NewResult(archiveisLabel, str, nil)) {
					resp.Body.Close()
					return
				}

			}

			resp.Body.Close()

			err = scanner.Err()

			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(archiveisLabel, nil, err))
				return
			}
		}

	}(domain, results)
	return results
}
