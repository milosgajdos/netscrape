package entity

// TypeFromString returns Type from string
// It returns error if the type can't be decoded from string.
func TypeFromString(s string) (Type, error) {
	switch s {
	case EntityString, "entity":
		return EntityType, nil
	case ResourceString, "resource":
		return ResourceType, nil
	}
	return UnknownType, ErrUnknownType
}
