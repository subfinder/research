package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/subfinder/research/core"
	"github.com/subfinder/research/core/sources"
)

// HTTPSSources is a list of all of the sources served over HTTPS.
var HTTPSSources = []core.Source{
	&sources.ArchiveIs{},
	&sources.CertSpotter{},
	&sources.CommonCrawlDotOrg{},
	&sources.CrtSh{},
	&sources.FindSubdomainsDotCom{},
	&sources.HackerTarget{},
	&sources.Riddler{},
	&sources.Threatminer{},
	&sources.WaybackArchive{},
	&sources.DNSDbDotCom{},
	&sources.Bing{},
	&sources.Yahoo{},
	&sources.Baidu{},
	&sources.Entrust{},
	&sources.ThreatCrowd{},
}

// HTTPSources is a list of all of the sources served over HTTP (plaintext).
var HTTPSources = []core.Source{
	&sources.PTRArchiveDotCom{},
	&sources.DogPile{},
}

// EnumerateOpts is used for the Enumerate function.
type EnumerateOpts struct {
	Domain          string `json:"domain"`
	IncludeFailures bool   `json:"include_failures"`
	AllSources      bool   `json:"all_sources"`
	HTTPSSources    bool   `json:"http_sources"`
	HTTPSources     bool   `json:"https_sources"`
}

// Enumerate is used for the /api/v1/enumerate endpoint.
func Enumerate(w http.ResponseWriter, r *http.Request) {
	uniqFilter := map[string]bool{}
	options := EnumerateOpts{}
	err := json.NewDecoder(r.Body).Decode(&options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if options.Domain == "" {
		fmt.Fprintf(w, "no domain given")
		return
	}
	enumOpts := core.EnumerationOptions{Sources: []core.Source{}}
	if options.AllSources {
		for _, source := range HTTPSSources {
			enumOpts.Sources = append(enumOpts.Sources, source)
		}
		for _, source := range HTTPSources {
			enumOpts.Sources = append(enumOpts.Sources, source)
		}
	}
	if options.HTTPSSources {
		for _, source := range HTTPSSources {
			enumOpts.Sources = append(enumOpts.Sources, source)
		}
	}
	if options.HTTPSources {
		for _, source := range HTTPSources {
			enumOpts.Sources = append(enumOpts.Sources, source)
		}
	}
	if len(enumOpts.Sources) == 0 {
		// default HTTPSSources
		for _, source := range HTTPSSources {
			enumOpts.Sources = append(enumOpts.Sources, source)
		}
	}
	for result := range core.EnumerateSubdomains(options.Domain, &enumOpts) {
		if result.IsSuccess() {
			_, found := uniqFilter[result.Success.(string)]
			if !found {
				uniqFilter[result.Success.(string)] = true
				if json, err := result.JSON(); err == nil {
					fmt.Fprintln(w, string(json))
				}
			}
		} else if options.IncludeFailures && result.IsFailure() {
			if json, err := result.JSON(); err == nil {
				fmt.Fprintln(w, string(json))
			}
		}
	}
}

func main() {
	http.HandleFunc("/api/v1/enumerate", Enumerate)

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
