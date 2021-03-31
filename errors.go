package netscrape

import "errors"

var (
	// ErrNotImplemented is returned when requesting a feature that has not been implemented yet.
	ErrNotImplemented = errors.New("ErrNotImplemented")
	// ErrMissingPlan is returned when no plan has been provided for netscraping.
	ErrMissingPlan = errors.New("ErrMissingPlan")
)
