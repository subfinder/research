package subzero

import "sync"
import "bytes"
import "errors"
import "fmt"
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
		Type:      t,
		Timestamp: time.Now().UTC(),
		Success:   s,
		Failure:   f,
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

// Printable turns a Result's information into a printable format (for STDOUT).
func (r *Result) Printable() string {
	var buffer bytes.Buffer
	r.RLock()
	defer r.RUnlock()

	if r.HasTimestamp() {
		buffer.WriteString(fmt.Sprintf("%v", r.Timestamp))
	}

	if r.HasType() {
		buffer.WriteString(fmt.Sprintf(" Type: %v", r.Type))
	}

	if r.IsSuccess() {
		buffer.WriteString(fmt.Sprintf(" Success: %v", r.Success))
	} else {
		buffer.WriteString(fmt.Sprintf(" Failure: %v", r.Failure))
	}

	return buffer.String()
}

// IsPrintable checks if the underlying Result has any printable information.
func (r *Result) IsPrintable() (bool, string) {
	printable := r.Printable()
	if len(printable) > 0 {
		return true, printable
	} else {
		return false, ""
	}
}

// Print will print the Printable version of the Result to the screen or return an error
// if the underlying Result has any printable information. Useful for debugging.
func (r *Result) Print() error {
	ok, printable := r.IsPrintable()
	if ok {
		fmt.Println(printable)
		return nil
	} else {
		return errors.New("unable to print unprintable result")
	}
}
