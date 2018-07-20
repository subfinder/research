package core

// Source defines the minimum interface any
// subdomain enumeration module should follow.
type Source interface {
	ProcessDomain(string) <-chan *Result
}
