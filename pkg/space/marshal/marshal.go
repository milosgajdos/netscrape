package marshal

import (
	"encoding/json"

	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
)

// Format defines encoding data format.
type Format int

const (
	// JSON data encoding.
	JSON Format = iota
)

// Marshal marshals m into format f.
func Marshal(f Format, m interface{}) ([]byte, error) {
	switch f {
	case JSON:
		return marshalJSON(m)
	default:
		return nil, ErrUnsupportedFormat
	}
}

// marshalJSON marshals m to JSON.
func marshalJSON(m interface{}) ([]byte, error) {
	switch v := m.(type) {
	case space.Entity:
		e := &Entity{
			UID:       v.UID().String(),
			Type:      v.Type(),
			Name:      v.Name(),
			Namespace: v.Namespace(),
			Attrs:     attrs.ToMap(v.Attrs()),
		}
		if v.Resource() != nil {
			e.Resource = &Resource{
				UID:        v.Resource().UID().String(),
				Type:       v.Resource().Type(),
				Name:       v.Resource().Name(),
				Group:      v.Resource().Group(),
				Version:    v.Resource().Version(),
				Kind:       v.Resource().Kind(),
				Namespaced: v.Resource().Namespaced(),
				Attrs:      attrs.ToMap(v.Resource().Attrs()),
			}
		}
		return json.Marshal(e)
	case space.Resource:
		r := &Resource{
			UID:        v.UID().String(),
			Type:       v.Type(),
			Name:       v.Name(),
			Group:      v.Group(),
			Version:    v.Version(),
			Kind:       v.Kind(),
			Namespaced: v.Namespaced(),
			Attrs:      attrs.ToMap(v.Attrs()),
		}
		return json.Marshal(r)
	case space.Link:
		l := &Link{
			UID:   v.UID().String(),
			From:  v.From().String(),
			To:    v.To().String(),
			Attrs: attrs.ToMap(v.Attrs()),
		}
		return json.Marshal(l)
	default:
		return nil, ErrUnsuportedType
	}
}

// Unmarshal unmarshals b in format f into m.
func Unmarshal(f Format, b []byte, m interface{}) error {
	switch f {
	case JSON:
		return unmarshalJSON(b, m)
	default:
		return ErrUnsupportedFormat
	}
}

// unmarshalJSON decodes JSON data stored in b into m.
func unmarshalJSON(b []byte, m interface{}) error {
	switch m.(type) {
	case *Entity, *Resource, *Link, *LinkedEntity:
		return json.Unmarshal(b, m)
	default:
		return ErrUnsuportedType
	}
}
