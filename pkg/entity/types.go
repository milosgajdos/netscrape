package entity

import (
	"bytes"
	"encoding/json"
)

// Resource is an entity resource.
type Resource struct {
	UID        string            `json:"uid"`
	Type       Type              `type:"type"`
	Name       string            `json:"name"`
	Group      string            `json:"group"`
	Version    string            `json:"version"`
	Kind       string            `json:"kind"`
	Namespaced bool              `json:"namespaced"`
	Attrs      map[string]string `json:"attrs,omitempty"`
}

// Link between two entities.
type Link struct {
	UID   string            `json:"uid"`
	From  string            `json:"from"`
	To    string            `json:"to"`
	Attrs map[string]string `json:"attrs,omitempty"`
}

// Entity is an arbitrary entity.
type Entity struct {
	UID       string            `json:"uid"`
	Type      Type              `type:"type"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Resource  *Resource         `json:"resource,omitempty"`
	Attrs     map[string]string `json:"attrs,omitempty"`
}

// LinkedEntity is Entity that has Links to other entities.
type LinkedEntity struct {
	Entity
	Links []Link `json:"links,omitempty"`
}

// Type is entity type
type Type int

const (
	EntityType Type = iota
	ResourceType
	UnknownType
)

const (
	EntityString   = "Entity"
	ResourceString = "Resource"
	UnknownString  = "Unknown"
)

// String implements fmt.Stringer
func (t Type) String() string {
	switch t {
	case EntityType:
		return EntityString
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
