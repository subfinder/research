package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "strings"
import "encoding/json"

type CertSpotter struct{}

type certspotterObject struct {
	DNSNames []string `json:"dns_names"`
}

func (source *CertSpotter) ProcessDomain(domain string) <-chan *core.Result {
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

		// get response from the API, optionally with an API key
		resp, err := httpClient.Get("https://certspotter.com/api/v0/certs?domain=" + domain)
		if err != nil {
			results <- &core.Result{Type: "certspotter", Failure: err}
			return
		}
		defer resp.Body.Close()

		certspotterData := []*certspotterObject{}
		err = json.NewDecoder(resp.Body).Decode(&certspotterData)
		if err != nil {
			results <- &core.Result{Type: "certspotter", Failure: err}
			return
		}
		for _, block := range certspotterData {
			for _, dnsName := range block.DNSNames {
				// Fix Wildcard subdomains containing asterisk before them
				if strings.Contains(dnsName, "*.") {
					dnsName = strings.Split(dnsName, "*.")[1]
				}
				results <- &core.Result{Type: "certspotter", Success: dnsName}
			}
		}
	}(domain, results)
	return results
}
