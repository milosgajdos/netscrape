package uuid

import "testing"

func TestNewFromString(t *testing.T) {
	s := "randomUID"

	uid, err := NewFromString(s)
	if err != nil {
		t.Fatalf("failed to create new uid from string %q: %v", s, err)
	}

	if s != uid.Value() {
		t.Errorf("expected: %s, got: %s", s, uid.Value())
	}
}

func TestNew(t *testing.T) {
	u1, err := New()
	if err != nil {
		t.Fatalf("failed to create new uid from string: %v", err)
	}

	u2, err := New()
	if err != nil {
		t.Fatalf("failed to create new uid from string: %v", err)
	}

	if len(u1.Value()) == 0 || len(u2.Value()) == 0 {
		t.Fatalf("empty uid returned")
	}

	if u1.Value() == u2.Value() {
		t.Errorf("non-unique uids generated")
	}
}
