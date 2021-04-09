package json

import (
	"encoding/json"

	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/marshal"
)

// Marshaler is JSON Marshaler.
type Marshaler struct{}

// NewMarshaler creates a new JSON marshaler and returns it.
func NewMarshaler(opts ...Option) (*Marshaler, error) {
	return &Marshaler{}, nil
}

// Marshal marshals x into JSON encoded bytes.
func (m *Marshaler) Marshal(x interface{}) ([]byte, error) {
	switch v := x.(type) {
	case space.Resource:
		r, err := marshal.ResourceFromSpace(v)
		if err != nil {
			return nil, err
		}
		return json.Marshal(r)
	case space.Object:
		o, err := marshal.ObjectFromSpace(v)
		if err != nil {
			return nil, err
		}
		return json.Marshal(o)
	case space.Entity:
		e, err := marshal.EntityFromSpace(v)
		if err != nil {
			return nil, err
		}
		return json.Marshal(e)
	case space.Link:
		l, err := marshal.LinkFromSpace(v)
		if err != nil {
			return nil, err
		}
		return json.Marshal(l)
	default:
		return nil, marshal.ErrUnsuportedType
	}
}

// Unmarshal unmarshals b to object o.
func (m *Marshaler) Unmarshal(b []byte, x interface{}) error {
	switch x := x.(type) {
	case *space.Resource:
		var r marshal.Resource
		if err := json.Unmarshal(b, &r); err != nil {
			return err
		}
		var err error
		*x, err = marshal.ResourceToSpace(r)
		if err != nil {
			return err
		}
	case *space.Object:
		var o marshal.Object
		if err := json.Unmarshal(b, &o); err != nil {
			return err
		}
		var err error
		*x, err = marshal.ObjectToSpace(o)
		if err != nil {
			return err
		}
	case *space.Entity:
		var e marshal.Entity
		if err := json.Unmarshal(b, &e); err != nil {
			return err
		}
		var err error
		*x, err = marshal.EntityToSpace(e)
		if err != nil {
			return err
		}
	case *space.Link:
		var l marshal.Link
		if err := json.Unmarshal(b, &l); err != nil {
			return err
		}
		var err error
		*x, err = marshal.LinkToSpace(l)
		if err != nil {
			return err
		}
	default:
		return marshal.ErrUnsuportedType
	}
	return nil
}
