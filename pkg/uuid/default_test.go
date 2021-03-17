package uuid

import "testing"

func TestNewFromString(t *testing.T) {
	s := "randomUID"

	uid, err := NewFromString(s)
	if err != nil {
		t.Fatalf("failed to create new uid from string %q: %v", s, err)
	}

	if s != uid.String() {
		t.Errorf("expected: %s, got: %s", s, uid.String())
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

	if len(u1.String()) == 0 || len(u2.String()) == 0 {
		t.Fatalf("empty uid returned")
	}

	if u1.String() == u2.String() {
		t.Errorf("non-unique uids generated")
	}
}
