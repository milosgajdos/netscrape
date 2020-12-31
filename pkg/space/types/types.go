package types

// Resource is Space resource
type Resource struct {
	Name       string                 `json:"name"`
	Group      string                 `json:"group"`
	Version    string                 `json:"version"`
	Kind       string                 `json:"kind"`
	Namespaced bool                   `json:"namespaced"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Link is a link between Space objects.
type Link struct {
	UID      string                 `json:"uid"`
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Object is a Space object
type Object struct {
	UID       string                 `json:"uid"`
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Resource  Resource               `json:"resource"`
	Links     []Link                 `json:"links"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
