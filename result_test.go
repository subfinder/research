package subzero

import "testing"
import "fmt"
import "errors"
import "time"
import "reflect"

func TestResult(t *testing.T) {
	var units = []struct {
		got *Result
		exp string
	}{
		{&Result{Type: "example", Success: "info.bing.com"}, "info.bing.com"},
		{&Result{Type: "example", Failure: errors.New("failed")}, "failed"},
	}
	for _, u := range units {
		if u.got.Failure != nil {
			if !reflect.DeepEqual(u.got.Failure.Error(), u.exp) {
				t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
			}
		} else {
			if !reflect.DeepEqual(u.got.Success, u.exp) {
				t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
			}
		}
	}
}

func TestNewResult(t *testing.T) {
	var units = []struct {
		got *Result
		exp *Result
	}{
		{NewResult("example", "info.bing.com", nil), &Result{Type: "example", Success: "info.bing.com"}},
		{NewResult("example", nil, errors.New("failed")), &Result{Type: "example", Failure: errors.New("failed")}},
	}
	for _, u := range units {
		u.got.Timestamp = time.Time{} // ensure this isn't the reason for failure
		if !reflect.DeepEqual(u.exp, u.got) {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResult_IsSuccess(t *testing.T) {
	var units = []struct {
		got *Result
		exp bool
	}{
		{&Result{Type: "example", Success: "info.bing.com"}, true},
		{&Result{Type: "example", Failure: errors.New("failed")}, false},
	}
	for _, u := range units {
		if u.got.IsSuccess() != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResult_IsFailure(t *testing.T) {
	var units = []struct {
		got *Result
		exp bool
	}{
		{&Result{Type: "example", Success: "info.bing.com"}, false},
		{&Result{Type: "example", Failure: errors.New("failed")}, true},
	}
	for _, u := range units {
		if u.got.IsFailure() != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResult_HasType(t *testing.T) {
	var units = []struct {
		got *Result
		exp bool
	}{
		{&Result{Type: "example", Success: "info.bing.com"}, true},
		{&Result{Type: "example", Failure: errors.New("failed")}, true},
		{&Result{}, false},
	}
	for _, u := range units {
		if u.got.HasType() != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResult_HasTimestamp(t *testing.T) {
	var units = []struct {
		got *Result
		exp bool
	}{
		{&Result{}, false},
		{NewResult("", "", nil), true},
	}
	for _, u := range units {
		if u.got.HasTimestamp() != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResult_Printable(t *testing.T) {
	var units = []struct {
		got *Result
		exp string
	}{
		{&Result{}, ""},
		{NewResult("", "", nil), "Success:"},
		{NewResult("example", "", nil), "Type: example Success:"},
		{NewResult("example", "a.b.com", nil), "Type: example Success: a.b.com"},
	}
	for _, u := range units {
		u.got.Timestamp = time.Time{} // ensure this isn't the reason for failure
		if u.got.Printable() != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got.Printable())
		}
	}
}

func TestResult_IsPrintable(t *testing.T) {
	var units = []struct {
		got *Result
		exp bool
	}{
		{&Result{}, false},
		{NewResult("", "", nil), true},
		{NewResult("example", "", nil), true},
		{NewResult("example", "a.b.com", nil), true},
	}
	for _, u := range units {
		u.got.Timestamp = time.Time{} // ensure this isn't the reason for failure
		if ok, _ := u.got.IsPrintable(); ok != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got.Printable())
		}
	}
}

// TODO: func TestResult_Print(t *testing.T) {}

func ExampleResult() {
	result := Result{Type: "example", Success: "info.bing.com"}
	if result.Failure != nil {
		fmt.Println(result.Type, ":", result.Failure)
	} else {
		fmt.Println(result.Type, ":", result.Success)
	}
	// Output: example : info.bing.com
}

func ExampleNewResult() {
	result := NewResult("example", "info.google.com", nil)

	if result.IsFailure() {
		fmt.Println(result.Failure.Error())
	} else {
		if result.Success.(string) == "info.google.com" {
			fmt.Println("found example in success")
		}
	}
	// Output: found example in success
}

func ExampleResult_IsSuccess() {
	result := Result{Success: "wiggle.github.com"}
	if result.IsSuccess() {
		fmt.Println(result.Success)
	}
	// Output: wiggle.github.com
}

func ExampleResult_IsFailure() {
	result := Result{Failure: errors.New("failed to party")}
	if result.IsFailure() {
		fmt.Println(result.Failure.Error())
	}
	// Output: failed to party
}

func ExampleResult_HasType() {
	result := Result{Type: "example"}
	fmt.Println(result.HasType())
	// Output: true
}

func ExampleResult_HasTimestamp() {
	result := Result{} // no Timestamp set
	fmt.Println(result.HasTimestamp())
	// Output: false
}

func ExampleResult_Printable() {
	result := NewResult("example", "ex.ample.com", nil)
	result.Timestamp = time.Time{} // set default timestamp
	printable := result.Printable()
	fmt.Println(printable)
	// Output: Type: example Success: ex.ample.com
}

func ExampleResult_IsPrintable() {
	result := NewResult("example", "ex.ample.com", nil)
	ok, _ := result.IsPrintable()
	fmt.Println(ok)
	// Output: true
}

func ExampleResult_Print() {
	result := NewResult("example", "ex.ample.com", nil)
	result.Timestamp = time.Time{} // set default timestamp
	result.Print()
	// Output: Type: example Success: ex.ample.com
}

func ExampleResult_JSON() {
	result := NewResult("example", "ex.ample.com", nil)
	result.Timestamp = time.Time{} // set default timestamp
	bytes, _ := result.JSON()
	fmt.Println(string(bytes))
	// Output: {"Timestamp":"0001-01-01T00:00:00Z","Type":"example","Success":"ex.ample.com","Failure":null}
}
