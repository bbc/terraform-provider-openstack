package containers

// CreateOptsV131 is a version of CreateOpts supporting the experimental registry parameter
type CreateOptsV131 struct {
	CreateOpts
	Registry string `json:"registry,omitempty"`
}
