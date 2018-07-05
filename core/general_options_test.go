package subzero

import "testing"
import "reflect"
import "time"
import "runtime"

func TestGeneralOptions(t *testing.T) {
	opts := &GeneralOptions{}
	var units = []struct {
		got interface{}
		exp interface{}
	}{
		{opts.Verbose, false},
		{opts.ColorSupport, false},
		{opts.AvailableCores, 0},
		{opts.Recursive, false},
		{opts.PassiveOnly, false},
		{opts.IgnoreErrors, false},
		{opts.OutputType, ""},
		{opts.OutputDir, ""},
		{len(opts.TargetDomains), 0},
		{len(opts.Sources), 0},
		{len(opts.Resolvers), 0},
	}
	for _, u := range units {
		if !reflect.DeepEqual(u.exp, u.got) {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

func TestDefaultDNSResolvers(t *testing.T) {
	var units = []struct {
		got interface{}
		exp interface{}
	}{
		{len(defaultDNSResolvers), 8},
	}
	for _, u := range units {
		if !reflect.DeepEqual(u.exp, u.got) {
			t.Fatalf("expected '%v', got '%v'", u.exp, u.got)
		}
	}
}

