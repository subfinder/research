package subzero

// Result contains the information from any given
// source. It's the Source author's job to set the
// type when returning a result. Upon success, a
// Source source should provide a string as the found
// subdomain. Upon Failure, the source should provide an error.
type Result struct {
	Type    string
	Success interface{}
	Failure error
}

// IsSuccess checks if the Result has any failure before
// determining if the result succeeded.
func (r *Result) IsSuccess() bool {
	if r.Failure != nil {
		return false
	}
	return true
}

// IsFailure checks if the Result has any failure before
// determining if the result failed.
func (r *Result) IsFailure() bool {
	if r.Failure != nil {
		return true
	}
	return false
}

// HasType checks if the Result has a type value set.
func (r *Result) HasType() bool {
	if r.Type != "" {
		return true
	}
	return false
}
