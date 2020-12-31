package uuid

import (
	"github.com/google/uuid"
)

// uid implements UID.
type uid struct {
	id string
}

// NewFromString returns new UID created from uid.
func NewFromString(s string) (uid, error) {
	return uid{
		id: s,
	}, nil
}

// New creates new UID and returns it.
func New() (uid, error) {
	return uid{
		id: uuid.New().String(),
	}, nil
}

// Value returns UID as a string
func (u uid) Value() string {
	return u.id
}
