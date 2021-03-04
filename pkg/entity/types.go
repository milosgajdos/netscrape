package entity

import (
	"bytes"
	"encoding/json"
)

// Type is entity type
type Type int

const (
	ObjectType Type = iota
	ResourceType
	UnknownType
)

const (
	ObjectString   = "Object"
	ResourceString = "Resource"
	UnknownString  = "Unknown"
)

// String implements fmt.Stringer
func (t Type) String() string {
	switch t {
	case ObjectType:
		return ObjectString
	case ResourceType:
		return ResourceString
	default:
		return UnknownString
	}
}

// MarshalJSON marshals t into JSOn encoded bytes.
func (t *Type) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(t.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals b into Type.
func (t *Type) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	typ, err := TypeFromString(s)
	if err != nil {
		return err
	}

	*t = typ
	return nil
}
