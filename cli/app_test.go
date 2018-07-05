package cli

import "testing"

func TestName(t *testing.T) {
	if Name != "SubFinder" {
		t.Error("expected cli application name to be 'SubFinder'")
	}
}

func TestVersion(t *testing.T) {
	if Version != "v2.0.0" {
		t.Error("expected cli version to be 'v2.0.0'")
	}
}

func TestNewApplication(t *testing.T) {
	app := NewApplication()
	if app.Name != "SubFinder" {
		t.Error("expected cli application name to be 'SubFinder'")
	}
	if app.Version != "v2.0.0" {
		t.Error("expected cli version to be 'v2.0.0'")
	}
	//app.Run(os.Args)
}
