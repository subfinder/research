package subzero

import "testing"
import "fmt"
import "errors"
import "time"
import "reflect"
import "sync"

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
		{&Result{Success: nil}, ""},
		{&Result{Failure: nil}, ""},
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
		{&Result{Success: nil}, false},
		{&Result{Failure: nil}, false},
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

func TestResult_JSON(t *testing.T) {
	var units = []struct {
		got *Result
		exp string
	}{
		{&Result{}, `{"Timestamp":"0001-01-01T00:00:00Z","Type":"","Success":null,"Failure":null}`},
		{&Result{Success: nil}, `{"Timestamp":"0001-01-01T00:00:00Z","Type":"","Success":null,"Failure":null}`},
		{&Result{Failure: nil}, `{"Timestamp":"0001-01-01T00:00:00Z","Type":"","Success":null,"Failure":null}`},
		{NewResult("", "", nil), `{"Timestamp":"0001-01-01T00:00:00Z","Type":"","Success":"","Failure":null}`},
		{NewResult("example", "", nil), `{"Timestamp":"0001-01-01T00:00:00Z","Type":"example","Success":"","Failure":null}`},
		{NewResult("example", "a.b.com", nil), `{"Timestamp":"0001-01-01T00:00:00Z","Type":"example","Success":"a.b.com","Failure":null}`},
	}
	for _, u := range units {
		u.got.Timestamp = time.Time{} // ensure this isn't the reason for failure
		if bytes, _ := u.got.JSON(); string(bytes) != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, string(bytes))
		}
	}
}

func TestResult_SetSuccess(t *testing.T) {
	var units = []struct {
		got *Result
		exp interface{}
	}{
		{&Result{}, ""},
		{&Result{Success: "cats"}, "dogs"},
		{&Result{Success: nil}, "birds"},
	}
	for _, u := range units {
		u.got.Timestamp = time.Time{} // ensure this isn't the reason for failure
		u.got.SetSuccess(u.exp)
		if u.got.Success != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
		if !u.got.IsSuccess() {
			t.Fatalf("expected '%v', to be a success", u.exp)
		}
	}
}

func TestResult_GetSuccess(t *testing.T) {
	var units = []struct {
		got *Result
		exp interface{}
	}{
		{&Result{}, ""},
		{&Result{Success: "cats"}, "dogs"},
		{&Result{Success: nil}, "birds"},
	}
	for _, u := range units {
		u.got.Timestamp = time.Time{} // ensure this isn't the reason for failure
		u.got.SetSuccess(u.exp)
		if u.got.GetSuccess() != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
		if !u.got.IsSuccess() {
			t.Fatalf("expected '%v', to be a success", u.exp)
		}
	}
}

func TestResult_GetType(t *testing.T) {
	var units = []struct {
		got *Result
		exp string
	}{
		{NewResult("lol", "", nil), "lol"},
		{NewResult("", "", nil), ""},
		{&Result{}, ""},
		{&Result{Type: "cats"}, "cats"},
		{&Result{Type: ""}, ""},
	}
	for _, u := range units {
		if u.got.GetType() != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResult_SetType(t *testing.T) {
	var units = []struct {
		got *Result
		exp string
	}{
		{NewResult("lol", "", nil), "lol"},
		{NewResult("", "", nil), ""},
		{NewResult("", "", nil), "lol"},
		{&Result{}, ""},
		{&Result{Type: "cats"}, "cats"},
		{&Result{Type: "dogs"}, "cats"},
		{&Result{Type: ""}, ""},
	}
	for _, u := range units {
		u.got.SetType(u.exp)
		if u.got.Type != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestResult_GetFailure(t *testing.T) {
	ex := errors.New("oh man")
	var units = []struct {
		got *Result
		exp error
	}{
		{NewResult("lol", "", ex), ex},
		{NewResult("lol", "", nil), nil},
	}
	for _, u := range units {
		if u.got.Failure != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got.Failure.Error())
		}
	}
}

func TestResult_GetTimestamp(t *testing.T) {
	var units = []struct {
		got *Result
		exp time.Time
	}{
		{NewResult("", "", nil), time.Time{}},
		{NewResult("", "", nil), time.Now()},
	}
	for _, u := range units {
		u.got.Timestamp = u.exp
		if u.got.GetTimestamp() != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got.Timestamp)
		}
	}
}

func TestResult_SetTimestamp(t *testing.T) {
	var units = []struct {
		got *Result
		exp time.Time
	}{
		{NewResult("", "", nil), time.Time{}},
		{NewResult("", "", nil), time.Now()},
	}
	for _, u := range units {
		u.got.SetTimestamp(u.exp)
		if u.got.Timestamp != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got.Timestamp)
		}
	}
}

func TestResult_SetFailure(t *testing.T) {
	ex := errors.New("oh man")
	var units = []struct {
		got *Result
		exp error
	}{
		{NewResult("lol", "", nil), ex},
		{NewResult("lol", "", nil), nil},
		{NewResult("lol", "", ex), errors.New("new error")},
		{NewResult("lol", "", ex), nil},
		{NewResult("lol", "", errors.New("anouther one!")), nil},
	}
	for _, u := range units {
		u.got.SetFailure(u.exp)
		if u.got.Failure != u.exp {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got.Failure)
		}
	}
}

func TestResultMultiThreadedBehavior(t *testing.T) {
	times := []struct {
		timeout int
		value   string
	}{
		{5, "a"},
		{3, "b"},
		{4, "c"},
		{4, "d"},
		{2, "e"},
		{3, "f"},
	}
	expWinner := "e" // smallest timeout

	sharedResult := &Result{}

	wg := sync.WaitGroup{}

	for _, t := range times {
		wg.Add(1)
		go func(t int, v string, r *Result) {
			defer wg.Done()
			// Note: this is only ok for Milliseconds, not
			// nanoseconds in my testing.
			time.Sleep(time.Duration(t) * time.Millisecond)
			if !r.IsSuccess() {
				r.SetSuccess(v)
			}
		}(t.timeout, t.value, sharedResult)
	}

	wg.Wait()

	if sharedResult.Success != expWinner {
		t.Fatalf("expected '%v', got '%v'", expWinner, sharedResult.Success)
	}
}

func TestResultMultiThreadedBehaviorMore(t *testing.T) {
	times := []struct {
		timeout int
		value   string
	}{
		{3, "a"},
		{4, "b"},
		{5, "c"},
		{4, "d"},
		{2, "e"},
		{3, "f"},
	}

	expWinner := "c" // last one to obtain a lock

	sharedResult := &Result{}

	wg := sync.WaitGroup{}

	for _, t := range times {
		wg.Add(1)
		go func(t int, v string, r *Result) {
			defer wg.Done()
			// Note: this is showing how the locks work when used in
			// a different, but still expected manner.
			if !r.IsSuccess() {
				time.Sleep(time.Duration(t) * time.Millisecond)
				r.SetSuccess(v)
			}
		}(t.timeout, t.value, sharedResult)
	}

	wg.Wait()

	if sharedResult.Success != expWinner {
		t.Fatalf("expected '%v', got '%v'", expWinner, sharedResult.Success)
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

func ExampleResult_JSON() {
	result := NewResult("example", "ex.ample.com", nil)
	result.Timestamp = time.Time{} // set default timestamp
	bytes, _ := result.JSON()
	fmt.Println(string(bytes))
	// Output: {"Timestamp":"0001-01-01T00:00:00Z","Type":"example","Success":"ex.ample.com","Failure":null}
}

func ExampleResult_SetSuccess() {
	result := NewResult("example", "", nil)
	// do work, possibly in multiple go routines
	result.SetSuccess([]string{"a.com", "b.com", "c.com"})
	// check if success
	fmt.Println(result.IsSuccess())
	// Output: true
}

func ExampleResult_GetSuccess() {
	result := NewResult("bing", "info.bing.com", nil)
	s := result.GetSuccess()
	fmt.Println(s)
	// Output: info.bing.com
}

func ExampleResult_GetType() {
	result := NewResult("bing", "info.bing.com", nil)
	t := result.GetType()
	fmt.Println(t)
	// Output: bing
}

func ExampleResult_SetType() {
	result := NewResult("bing", "info.bing.com", nil)
	result.SetType("google")
	fmt.Println(result.Type)
	// Output: google
}

func ExampleResult_SetTimestamp() {
	result := NewResult("bing", "info.bing.com", nil)
	newTimestamp := time.Now().UTC()
	result.SetTimestamp(newTimestamp)
	fmt.Println(result.Timestamp == newTimestamp)
	// Output: true
}

func ExampleResult_GetTimestamp() {
	result := NewResult("bing", "info.bing.com", nil)
	newTimestamp := time.Now().UTC()
	result.Timestamp = newTimestamp
	fmt.Println(result.GetTimestamp() == newTimestamp)
	// Output: true
}

func ExampleResult_GetFailure() {
	McErr := errors.New("whoa there!")
	result := NewResult("", nil, McErr)
	fmt.Println(result.Failure)
	// Output: whoa there!
}

func ExampleResult_SetFailure() {
	McErr := errors.New("whoa there!")
	result := NewResult("thisis", "totally.fine.com", nil)
	result.SetFailure(McErr)
	if result.IsFailure() {
		fmt.Println(result.Failure, "we found our failure!")
	}
	if result.IsSuccess() {
		fmt.Println("this will never print because the failure was set")
	}
	// Output: whoa there! we found our failure!
}

func BenchmarkResultGetTypeThreadSafe(b *testing.B) {
	r := NewResult("example", "picat was here", nil)
	for n := 0; n < b.N; n++ {
		if r.GetType() == "example" {
		}
	}
}

func BenchmarkResultGetTypeThreadSafeMultiThreaded(b *testing.B) {
	r := NewResult("example", "picat was here", nil)
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func(r *Result) {
			defer wg.Done()
			if r.GetType() == "example" {
			}
		}(r)
	}
	wg.Wait()
}

func BenchmarkResultGetType(b *testing.B) {
	r := NewResult("example", "picat was here", nil)
	for n := 0; n < b.N; n++ {
		if r.Type == "example" {
			continue
		}
	}
}

func BenchmarkResultGetTypeMultiThreaded(b *testing.B) {
	r := NewResult("example", "picat was here", nil)
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func(r *Result) {
			defer wg.Done()
			if r.Type == "example" {
			}
		}(r)
	}
	wg.Wait()
}

func BenchmarkNewResultSingleThreaded(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewResult("example", n, nil)
	}
}

func BenchmarkNewResultMultiThreaded(b *testing.B) {
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			NewResult("example", n, nil)
		}()
	}
	wg.Wait()
}
