package sources

import (
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"

	core "github.com/subfinder/research/core"
)

type Ipv4Info struct{}

func (source *Ipv4Info) ProcessDomain(domain string) <-chan *core.Result {
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
			results <- &core.Result{Type: "ipv4info", Failure: err}
			return
		}
		regIPAddressPageToken := regexp.MustCompile("/ip-address/(.*)/" + domain)
		regxDNSAddressPageToken := regexp.MustCompile("/dns/(.*?)/" + domain)
		regxFirstSubdomainPageToken := regexp.MustCompile("/subdomains/(.*?)/" + domain)

		uniqFilter := map[string]bool{}

		body, err := getBody(httpClient, "http://ipv4info.com/search/"+domain)
		if err != nil {
			results <- &core.Result{Type: "ipv4info", Failure: err}
			return
		}

		// Get IP address page token
		matchTokens := regIPAddressPageToken.FindAllString(body, 1)
		if len(matchTokens) == 0 {
			results <- &core.Result{Type: "ipv4info", Failure: err}
			return
		}

		body, err = getBody(httpClient, "http://ipv4info.com"+matchTokens[0])
		if err != nil {
			results <- &core.Result{Type: "ipv4info", Failure: err}
			return
		}

		// Get DNS address page token
		matchTokens = regxDNSAddressPageToken.FindAllString(body, -1)
		if len(matchTokens) == 0 {
			results <- &core.Result{Type: "ipv4info", Failure: err}
			return
		}

		body, err = getBody(httpClient, "http://ipv4info.com"+matchTokens[0])
		if err != nil {
			results <- &core.Result{Type: "ipv4info", Failure: err}
			return
		}

		// Get First Subdomains page token
		matchTokens = regxFirstSubdomainPageToken.FindAllString(body, -1)
		if len(matchTokens) == 0 {
			results <- &core.Result{Type: "ipv4info", Failure: err}
			return
		}

		// Get first subdomains page
		body, err = getBody(httpClient, "http://ipv4info.com"+matchTokens[0])
		if err != nil {
			results <- &core.Result{Type: "ipv4info", Failure: err}
			return
		}

		nextPage := 1

		for {
			regxTokens := regexp.MustCompile("/subdomains/.*/page" + strconv.Itoa(nextPage) + "/" + domain + ".html")
			matchTokens := regxTokens.FindAllString(body, -1)
			if len(matchTokens) == 0 {
				return
			}

			body, err = getBody(httpClient, "http://ipv4info.com"+matchTokens[0])
			if err != nil {
				results <- &core.Result{Type: "ipv4info", Failure: err}
				return
			}

			for _, str := range domainExtractor.FindAllString(body, -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- &core.Result{Type: "ipv4info", Success: str}
				}
			}
			nextPage++
		}
	}(domain, results)
	return results
}

func getBody(httpClient *http.Client, url string) (string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
