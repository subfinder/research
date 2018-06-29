package subzero

import "sync"
import "time"

// Result contains the information from any given
// source. It's the Source author's job to set the
// type when returning a result. Upon success, a
// Source source should provide a string as the found
// subdomain. Upon Failure, the source should provide an error.
type Result struct {
	sync.RWMutex
	Timestamp time.Time
	Type      string
	Success   interface{}
	Failure   error
}

// NewResult wraps up the creation of a new Result.
func NewResult(t string, s interface{}, f error) *Result {
	return &Result{
		Type:    t,
		Success: s,
		Failure: f,
	}
}

// IsSuccess checks if the Result has any failure before
// determining if the result succeeded.
func (r *Result) IsSuccess() bool {
	r.RLock()
	defer r.RUnlock()
	if r.Failure != nil {
		return false
	}
	return true
}

// IsFailure checks if the Result has any failure before
// determining if the result failed.
func (r *Result) IsFailure() bool {
	r.RLock()
	defer r.RUnlock()
	if r.Failure != nil {
		return true
	}
	return false
}

// HasType checks if the Result has a type value set.
func (r *Result) HasType() bool {
	r.RLock()
	defer r.RUnlock()
	if r.Type != "" {
		return true
	}
	return false
}

// defaultTimestampValue is a cached variable used in HasTimestamp.
var defaultTimestampValue = time.Time{}

// HasTimestamp checks if the Result has a timestamp set.
func (r *Result) HasTimestamp() bool {
	r.RLock()
	defer r.RUnlock()
	if r.Timestamp != defaultTimestampValue {
		return true
	}
	return false
}
