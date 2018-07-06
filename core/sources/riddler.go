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

