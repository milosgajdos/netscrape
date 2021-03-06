package entity

import (
	"bytes"
	"testing"
)

func TestTypeString(t *testing.T) {
	testCases := []struct {
		t Type
		s string
	}{
		{EntityType, EntityString},
		{ResourceType, ResourceString},
		{UnknownType, UnknownString},
		{Type(-100), UnknownString},
	}

	for _, tc := range testCases {
		if s := tc.t.String(); s != tc.s {
			t.Errorf("expected string: %s, got: %s", tc.s, s)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	testCases := []struct {
		t   Type
		b   []byte
		err error
	}{
		{EntityType, []byte(`"` + EntityString + `"`), nil},
		{ResourceType, []byte(`"` + ResourceString + `"`), nil},
		{UnknownType, []byte(`"` + UnknownString + `"`), nil},
		{Type(-100), []byte(`"` + UnknownString + `"`), nil},
	}

	for _, tc := range testCases {
		b, err := tc.t.MarshalJSON()
		if !bytes.Equal(b, tc.b) {
			t.Errorf("expected bytes: %v, got;: %v", tc.b, b)
		}

		if err != tc.err {
			t.Errorf("expected error: %v, got: %v", tc.err, err)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		t   Type
		b   []byte
		err error
	}{
		{EntityType, []byte(`"` + EntityString + `"`), nil},
		{ResourceType, []byte(`"` + ResourceString + `"`), nil},
		{UnknownType, []byte(`"` + UnknownString + `"`), ErrUnknownType},
		{Type(-100), []byte(`"` + UnknownString + `"`), ErrUnknownType},
	}

	for _, tc := range testCases {
		if err := tc.t.UnmarshalJSON(tc.b); err != tc.err {
			t.Errorf("expected error: %v, got: %v", tc.err, err)
		}
	}
}
