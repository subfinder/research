package subzero

import "testing"
import "fmt"
import "errors"
import "time"
import "reflect"

func TestResult(t *testing.T) {
	var units = []struct {
		exp *Result
		got string
	}{
		{&Result{Type: "example", Success: "info.bing.com"}, "info.bing.com"},
		{&Result{Type: "example", Failure: errors.New("failed")}, "failed"},
	}
	for _, u := range units {
		if u.exp.Failure != nil {
			if !reflect.DeepEqual(u.exp.Failure.Error(), u.got) {
				t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
			}
		} else {
			if !reflect.DeepEqual(u.exp.Success, u.got) {
				t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
			}
		}
	}
}

func TestNewResult(t *testing.T) {
	var units = []struct {
		exp *Result
		got *Result
	}{
		{&Result{Type: "example", Success: "info.bing.com"}, NewResult("example", "info.bing.com", nil)},
		{&Result{Type: "example", Failure: errors.New("failed")}, NewResult("example", nil, errors.New("failed"))},
	}
	for _, u := range units {
		u.got.Timestamp = time.Time{} // ensure this isn't the reason for failure
		if !reflect.DeepEqual(u.exp, u.got) {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResultIsSuccess(t *testing.T) {
	var units = []struct {
		exp *Result
		got bool
	}{
		{&Result{Type: "example", Success: "info.bing.com"}, true},
		{&Result{Type: "example", Failure: errors.New("failed")}, false},
	}
	for _, u := range units {
		if u.exp.IsSuccess() != u.got {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResultIsFailure(t *testing.T) {
	var units = []struct {
		exp *Result
		got bool
	}{
		{&Result{Type: "example", Success: "info.bing.com"}, false},
		{&Result{Type: "example", Failure: errors.New("failed")}, true},
	}
	for _, u := range units {
		if u.exp.IsFailure() != u.got {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResultHasType(t *testing.T) {
	var units = []struct {
		exp *Result
		got bool
	}{
		{&Result{Type: "example", Success: "info.bing.com"}, true},
		{&Result{Type: "example", Failure: errors.New("failed")}, true},
		{&Result{}, false},
	}
	for _, u := range units {
		if u.exp.HasType() != u.got {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

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
