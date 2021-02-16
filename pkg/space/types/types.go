package types

// Resource is space resource.
type Resource struct {
	UID        string            `json:"uid"`
	Name       string            `json:"name"`
	Group      string            `json:"group"`
	Version    string            `json:"version"`
	Kind       string            `json:"kind"`
	Namespaced bool              `json:"namespaced"`
	Attrs      map[string]string `json:"attrs,omitempty"`
}

// Link between two space entities.
type Link struct {
	UID   string            `json:"uid"`
	From  string            `json:"from"`
	To    string            `json:"to"`
	Attrs map[string]string `json:"attrs,omitempty"`
}

// Object is space object
type Object struct {
	UID       string            `json:"uid"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Resource  Resource          `json:"resource"`
	Links     []Link            `json:"links"`
	Attrs     map[string]string `json:"attrs,omitempty"`
}
