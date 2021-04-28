package memory

import "testing"

func TestNewFromString(t *testing.T) {
	s := "randomUID"

	uid := NewFromString(s)

	if s != uid.String() {
		t.Errorf("expected: %s, got: %s", s, uid.String())
	}
}

func TestNew(t *testing.T) {
	u1 := New()
	if len(u1.String()) == 0 {
		t.Fatalf("empty uid created")
	}

	u2 := New()
	if len(u2.String()) == 0 {
		t.Fatalf("empty uid created")
	}

	if u1.String() == u2.String() {
		t.Errorf("non-unique uids generated")
	}
}
