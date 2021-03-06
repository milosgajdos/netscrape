package entity

import "testing"

func TestTypeFromString(t *testing.T) {
	testCases := []struct {
		s   string
		t   Type
		err error
	}{
		{"Resource", ResourceType, nil},
		{"resource", ResourceType, nil},
		{"entity", EntityType, nil},
		{"Entity", EntityType, nil},
		{"foo", UnknownType, ErrUnknownType},
	}

	for _, tc := range testCases {
		typ, err := TypeFromString(tc.s)
		if typ != tc.t {
			t.Errorf("expected type: %v, got: %v", tc.t, typ)
		}
		if err != tc.err {
			t.Errorf("expected error: %v, got: %v", tc.err, err)
		}
	}
}
