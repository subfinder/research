package sources

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/subfinder/research/core"
)

// Riddler is a source to process subdomains from https://riddler.io
type Riddler struct {
	Email    string
	Password string
	APIToken string
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
func (source *Riddler) Authenticate() (bool, error) {
	var data = []byte(`{"email":"` + source.Email + `", "password":"` + source.Password + `"}`)

	// Create a post request to get subdomain data
	req, err := http.NewRequest("POST", "https://riddler.io/auth/login", bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/json")

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
func (source *Riddler) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		// check if only email was given
		if source.Email != "" && source.Password == "" {
			results <- core.NewResult("riddler", nil, errors.New("given email, but no password"))
		}

		// check if only password was given
		if source.Email == "" && source.Password != "" {
			results <- core.NewResult("riddler", nil, errors.New("given password, but no email"))
		}

		// check if source needs to be authenticated
		if source.APIToken == "" && source.Email != "" && source.Password != "" {
			_, err := source.Authenticate()
			if err != nil {
				results <- core.NewResult("riddler", nil, err)
				return
			}
		}

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("riddler", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		var resp *http.Response

		if source.APIToken != "" {
			query := strings.NewReader(`{"query": "pld:` + domain + `", "output": "host", "limit": 500}`)
			req, err := http.NewRequest("POST", "https://riddler.io/api/search", query)
			if err != nil {
				// handle err
			}
			req.Header.Set("Content-type", "application/json")
			req.Header.Set("Authentication-Token", source.APIToken)

			resp, err := core.HTTPClient.Do(req)
			if err != nil {
				results <- core.NewResult("riddler", nil, err)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				results <- core.NewResult("riddler", nil, errors.New(resp.Status))
				return
			}

			hostResponse := []*riddlerHost{}

			err = json.NewDecoder(resp.Body).Decode(&hostResponse)
			if err != nil {
				results <- core.NewResult("riddler", nil, err)
				return
			}

			for _, r := range hostResponse {
				for _, str := range domainExtractor.FindAllString(r.Host, -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						results <- core.NewResult("riddler", str, nil)
					}
				}
			}
			return
		}

		if source.APIToken == "" {
			// not authenticated
			resp, err = core.HTTPClient.Get("https://riddler.io/search/exportcsv?q=pld:" + domain)
			if err != nil {
				results <- core.NewResult("riddler", nil, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				results <- core.NewResult("riddler", nil, errors.New(resp.Status))
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			for scanner.Scan() {
				for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						results <- core.NewResult("riddler", str, nil)
					}
				}
			}
			return
		}

	}(domain, results)
	return results
}
