package core

import (
	"regexp"
	"sync"
)

var subdomainExtractorMutex = &sync.Mutex{}

// NewSubdomainExtractor creates a new regular expression to extract
// subdomains from text based on the given domain.
func NewSubdomainExtractor(domain string) (*regexp.Regexp, error) {
	subdomainExtractorMutex.Lock()
	defer subdomainExtractorMutex.Unlock()
	extractor, err := regexp.Compile(`[\w-\*]+\.` + domain)
	if err != nil {
		return nil, err
	}
	return extractor, nil
}
