package subzero

import "testing"
import "fmt"

func TestEnumerationOptions(t *testing.T) {
	opts := EnumerationOptions{}
	if len(opts.Sources) != 0 {
		t.Error("new enumeration options doesn't have empty sources list")
	}
}

func TestEnumerationOptionsHasSources(t *testing.T) {
	opts := EnumerationOptions{}
	if opts.HasSources() {
		t.Error("new enumeration options doesn't have empty sources list")
	}
}

func ExampleEnumerationOptions() {
	opts := EnumerationOptions{}
	if opts.HasSources() {
		fmt.Println("sources found in options")
	} else {
		fmt.Println("sources not found in options")
	}
	// Output: sources not found in options
}

func ExampleEnumerationOptions_HasSources() {
	opts := EnumerationOptions{}
	fmt.Println(opts.HasSources())
	// Output: false
}
