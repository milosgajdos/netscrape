package marshal

// Entity is an arbitrary entity.
type Entity struct {
	UID   string            `json:"uid"`
	Type  string            `json:"type"`
	Attrs map[string]string `json:"attrs,omitempty"`
}

// Resource is an arbitrary resource.
type Resource struct {
	Entity
	Name       string `json:"name"`
	Group      string `json:"group"`
	Version    string `json:"version"`
	Kind       string `json:"kind"`
	Namespaced bool   `json:"namespaced"`
}

// Object is an arbitrary object.
type Object struct {
	Entity
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Resource  *Resource `json:"resource,omitempty"`
}

// Link between two entities.
type Link struct {
	UID   string            `json:"uid"`
	From  string            `json:"from"`
	To    string            `json:"to"`
	Attrs map[string]string `json:"attrs,omitempty"`
}

// LinkedObject is an Object linked to other objects.
type LinkedObject struct {
	Object
	Links []Link `json:"links,omitempty"`
}

// LinkedEntity is an Entity linked to other entities.
type LinkedEntity struct {
	Entity
	Links []Link `json:"links,omitempty"`
}
