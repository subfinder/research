package sources

import (
	"bufio"
	"strconv"

	"github.com/subfinder/research/core"
)

// DogPile is a source to process subdomains from http://dogpile.com
//
// Note
//
// This source uses http instead of https because of problems dogpile's SSL cert.
//
type DogPile struct{}

