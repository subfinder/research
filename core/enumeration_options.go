package core

import "time"

// EnumerationOptions provides all the data needed for subdomain
// enumeration. This includes all the sources which will be
// queried to find them.
type EnumerationOptions struct {
	Sources []Source
	Timeout time.Duration
}

// HasSources checks if the EnumerationOptions have any source defined.
func (opts *EnumerationOptions) HasSources() bool {
	if len(opts.Sources) == 0 {
		return false
	}
	return true
}
