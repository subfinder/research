package sources

import (
	"bufio"
	"errors"
	"strconv"

	"github.com/subfinder/research/core"
)

// Baidu is a source to process subdomains from https://baidu.com
type Baidu struct{}

