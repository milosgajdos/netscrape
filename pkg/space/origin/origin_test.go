package origin

import "testing"

const (
	src = "file://origin"
)

func TestNew(t *testing.T) {
	o, err := New(src)
	if err != nil {
		t.Fatalf("failed to create new origin: %v", err)
	}

	if u := o.URL().String(); u != src {
		t.Fatalf("expected origin URL: %s, got: %v", src, u)
	}

	if _, err = New(":foo"); err == nil {
		t.Fatal("expected error")
	}
}
