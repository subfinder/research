package sources

import (
	"bufio"
	"errors"
	"strconv"

	"github.com/subfinder/research/core"
)

// Ask is a source to process subdomains from https://ask.com
type Ask struct{}

