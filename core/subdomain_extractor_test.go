package subzero

import "testing"
import "fmt"

func TestNewSubdomainExtractor(t *testing.T) {
	_, err := NewSubdomainExtractor("google.com")
	if err != nil {
		t.Error(err)
	}
}

