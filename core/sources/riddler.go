package sources

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/subfinder/research/core"
	"golang.org/x/sync/semaphore"
)

// Riddler is a source to process subdomains from https://riddler.io
type Riddler struct {
	Email    string
	Password string
	APIToken string
	lock     *semaphore.Weighted
}

type riddlerHost struct {
	Host string `json:"host"`
}

type riddlerAuthenticationResponse struct {
	Response struct {
		User struct {
			AuthenticationToken string `json:"authentication_token"`
		} `json:"user"`
	} `json:"response"`
}

// Authenticate uses a given username and password to retrieve the APIToken.
func (source *Riddler) Authenticate(ctx context.Context) (bool, error) {
	var data = []byte(`{"email":"` + source.Email + `", "password":"` + source.Password + `"}`)

	// Create a post request to get subdomain data
	req, err := http.NewRequest("POST", "https://riddler.io/auth/login", bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/json")

	req.Cancel = ctx.Done()
	req.WithContext(ctx)

	resp, err := core.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	auth := &riddlerAuthenticationResponse{}

	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		return false, err
	}

	if auth.Response.User.AuthenticationToken == "" {
		return false, errors.New("failed to get authentication token")
	}

	source.APIToken = auth.Response.User.AuthenticationToken

	return true, nil
}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Riddler) ProcessDomain(ctx context.Context, domain string) <-chan *core.Result {
	if source.lock == nil {
		source.lock = defaultLockValue()
	}

	results := make(chan *core.Result)

	go func(domain string, results chan *core.Result) {
		defer close(results)

		if err := source.lock.Acquire(ctx, 1); err != nil {
			sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, err))
			return
		}
		defer source.lock.Release(1)

		// check if only email was given
		if source.Email != "" && source.Password == "" {
			sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, errors.New("given email, but no password")))
		}

		// check if only password was given
		if source.Email == "" && source.Password != "" {
			sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, errors.New("given password, but no email")))
		}

		// check if source needs to be authenticated
		if source.APIToken == "" && source.Email != "" && source.Password != "" {
			_, err := source.Authenticate(ctx)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, err))
				return
			}
		}

		domainExtractor := core.NewSingleSubdomainExtractor(domain)

		var resp *http.Response
		var err error

		if source.APIToken != "" {
			query := strings.NewReader(`{"query": "pld:` + domain + `", "output": "host", "limit": 500}`)
			req, err := http.NewRequest("POST", "https://riddler.io/api/search", query)
			if err != nil {
				// handle err
			}
			req.Header.Set("Content-type", "application/json")
			req.Header.Set("Authentication-Token", source.APIToken)

			req.WithContext(ctx)

			resp, err := core.HTTPClient.Do(req)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, err))
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, errors.New(resp.Status)))
				return
			}

			hostResponse := []*riddlerHost{}

			err = json.NewDecoder(resp.Body).Decode(&hostResponse)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, err))
				return
			}

			for _, r := range hostResponse {
				str := domainExtractor([]byte(r.Host))
				if str != "" {
					if !sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, str, nil)) {
						return
					}
				}
			}
			return
		}

		if source.APIToken == "" {
			// not authenticated
			resp, err = core.HTTPClient.Get("https://riddler.io/search/exportcsv?q=pld:" + domain)
			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, err))
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, errors.New(resp.Status)))
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			for scanner.Scan() {
				str := domainExtractor(scanner.Bytes())
				if str != "" {
					if !sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, str, nil)) {
						return
					}
				}
			}

			err = scanner.Err()

			if err != nil {
				sendResultWithContext(ctx, results, core.NewResult(riddlerLabel, nil, err))
				return
			}

			return
		}

	}(domain, results)
	return results
}
