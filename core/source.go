package core

import "context"

// Source defines the minimum interface any
// subdomain enumeration module should follow.
type Source interface {
	ProcessDomain(context.Context, string) <-chan *Result
}
