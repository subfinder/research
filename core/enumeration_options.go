package core

import "context"

// EnumerationOptions provides all the data needed for subdomain
// enumeration. This includes all the sources which will be
// queried to find them.
type EnumerationOptions struct {
	Sources   []Source
	Context   context.Context
	Cancel    context.CancelFunc
	Recursive bool
}

// HasSources checks if the EnumerationOptions have any source defined.
func (opts *EnumerationOptions) HasSources() bool {
	if len(opts.Sources) == 0 {
		return false
	}
	return true
}
