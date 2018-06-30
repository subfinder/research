package subzero

import "sync"
import "bytes"
import "errors"
import "fmt"
import "time"
import "strings"
import "encoding/json"

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

// SetSuccess safely sets a new success value to a Result
// which could be shared by multiple go routines.
func (r *Result) SetSuccess(success interface{}) {
	r.Lock()
	defer r.Unlock()
	r.Success = success
}

// GetSuccess safely gets a new success value from a Result
// which could be shared by multiple go routines.
func (r *Result) GetSuccess() interface{} {
	r.RLock()
	defer r.RUnlock()
	return r.Success
}

// SetFailure safely sets a new error value for a Result
// which could be shared by multiple go routines.
func (r *Result) SetFailure(err error) {
	r.Lock()
	defer r.Unlock()
	r.Failure = err
}

// GetFailure safely sets a new error value for a Result
// which could be shared by multiple go routines.
func (r *Result) GetFailure() error {
	r.RLock()
	defer r.RUnlock()
	return r.Failure
}

// SetType safely sets a new type string value for a Result
// which could be shared by multiple go routines.
func (r *Result) SetType(t string) {
	r.Lock()
	defer r.Unlock()
	r.Type = t
}

// IsSuccess checks if the Result has any failure before
// determining if the result succeeded.
func (r *Result) IsSuccess() bool {
	r.RLock()
	defer r.RUnlock()
	if r.Failure != nil {
		return false
	}
	// don't give us any of your bullshit, empty interfaces
	if fmt.Sprintf("%v", r.Success) == "<nil>" {
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
	}

	if r.IsFailure() {
		buffer.WriteString(fmt.Sprintf(" Failure: %v", r.Failure))
	}

	return strings.TrimSpace(buffer.String())
}

// IsPrintable checks if the underlying Result has any printable information.
func (r *Result) IsPrintable() (bool, string) {
	printable := r.Printable()
	if len(printable) > 0 {
		return true, printable
	}
	return false, ""
}

// Print will print the Printable version of the Result to the screen or return an error
// if the underlying Result has any printable information. Useful for debugging.
func (r *Result) Print() error {
	ok, printable := r.IsPrintable()
	if ok {
		fmt.Println(printable)
		return nil
	}
	return errors.New("unable to print unprintable result")
}

// JSON returns the Result as a JSON object within a slice of bytes.
func (r *Result) JSON() ([]byte, error) {
	r.RLock()
	defer r.RUnlock()
	return json.Marshal(r)
}
