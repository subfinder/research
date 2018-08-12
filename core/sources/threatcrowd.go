package sources

import (
	"bufio"
	"errors"

	"github.com/subfinder/research/core"
)

// ThreatCrowd is a source to process subdomains from https://threatcrowd.com
type ThreatCrowd struct{}

