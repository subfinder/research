package sources

import (
	"bufio"
	"errors"
	"strings"

	"github.com/subfinder/research/core"
)

// Entrust is a source to process subdomains from https://entrust.com
type Entrust struct{}

