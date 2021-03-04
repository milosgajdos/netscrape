package entity

type Entity struct {
	UID   string            `json:"uid"`
	Type  Type              `type:"type"`
	Attrs map[string]string `json:"attrs,omitempty"`
}

// Resource is resource entity
type Resource struct {
	Entity
	Name       string `json:"name"`
	Group      string `json:"group"`
	Version    string `json:"version"`
	Kind       string `json:"kind"`
	Namespaced bool   `json:"namespaced"`
}

// Link between two entities.
type Link struct {
	UID   string            `json:"uid"`
	From  string            `json:"from"`
	To    string            `json:"to"`
	Attrs map[string]string `json:"attrs,omitempty"`
}

// Object is object entity
type Object struct {
	Entity
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Resource  Resource `json:"resource"`
	Links     []Link   `json:"links"`
}
