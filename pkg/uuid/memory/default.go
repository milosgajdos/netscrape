package memory

import (
	"github.com/google/uuid"
)

// UID implements UID.
type UID struct {
	id string
}

// NewFromString returns a new UID created from UID.
func NewFromString(s string) *UID {
	return &UID{
		id: s,
	}
}

// New creates new UID and returns it.
func New() *UID {
	return &UID{
		id: uuid.New().String(),
	}
}

// String returns UID as a string
func (u UID) String() string {
	return u.id
}
