package marshal

// Resource is an entity resource.
type Resource struct {
	UID        string            `json:"uid"`
	Type       string            `type:"type"`
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
	Type      string            `type:"type"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Resource  *Resource         `json:"resource,omitempty"`
	Attrs     map[string]string `json:"attrs,omitempty"`
}

// LinkedEntity is Entity that is linked to other entities.
type LinkedEntity struct {
	Entity
	Links []Link `json:"links,omitempty"`
}
