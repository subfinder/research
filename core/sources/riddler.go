package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "strings"
import "encoding/json"
import "bytes"
import "bufio"
import "net"
import "time"
import "errors"

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

// Note: probably not totally thread safe? #yolo
func (source *Riddler) Authenticate() (bool, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	var data = []byte(`{"email":"` + source.Email + `", "password":"` + source.Password + `"}`)

	// Create a post request to get subdomain data
	req, err := http.NewRequest("POST", "https://riddler.io/auth/login", bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
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

func (source *Riddler) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		// check if source needs to be authenticated
		if source.APIToken == "" && source.Email != "" && source.Password != "" {
			_, err := source.Authenticate()
			if err != nil {
				results <- &core.Result{Type: "riddler", Failure: err}
				return
			}
		}

		httpClient := &http.Client{
			Timeout: time.Second * 60,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		}

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

			resp, err := httpClient.Do(req)
			if err != nil {
				results <- &core.Result{Type: "riddler", Failure: err}
				return
			}

			defer resp.Body.Close()

			hostResponse := []*riddlerHost{}

			err = json.NewDecoder(resp.Body).Decode(&hostResponse)
			if err != nil {
				results <- &core.Result{Type: "riddler", Failure: err}
				return
			}

			for _, r := range hostResponse {
				results <- &core.Result{Type: "riddler", Success: r.Host}
			}
			return
		}

		if source.APIToken == "" {
			// not authenticated
			resp, err = httpClient.Get("https://riddler.io/search/exportcsv?q=pld:" + domain)
			if err != nil {
				results <- &core.Result{Type: "riddler", Failure: err}
				return
			}
			defer resp.Body.Close()
			scanner := bufio.NewScanner(resp.Body)
			counter := 0
			for scanner.Scan() {
				if counter <= 2 {
					counter += 1
					continue // skip first two lines
				}
				strParts := strings.Split(scanner.Text(), ",")
				if (len(strParts) >= 6) && strings.Contains(strParts[5], domain) {
					results <- &core.Result{Type: "riddler", Success: strParts[5]}
				}
			}
			return
		}

	}(domain, results)
	return results
}
