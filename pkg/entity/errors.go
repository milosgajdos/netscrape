package entity

import "errors"

var (
	// ErrUnknownType is returned when decoding unknown tity type.
	ErrUnknownType = errors.New("ErrUnknownType")
)
