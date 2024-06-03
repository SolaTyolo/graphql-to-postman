package postman

type IRequest struct {
	Body        *Body    `json:"body,omitempty"`
	Description string   `json:"description,omitempty"`
	Header      []Header `json:"header,omitempty"`
	Method      string   `json:"method,omitempty"`
	URL         string   `json:"url,omitempty"`
}

// This field contains the data usually contained in the request body.
type Body struct {
	// When set to true, prevents request body from being sent.
	Graphql map[string]string `json:"graphql,omitempty"`
	// Postman stores the type of data associated with this request in this field.
	Mode *Mode `json:"mode,omitempty"`
	// Additional configurations and options set for various body modes.
	Options map[string]interface{} `json:"options,omitempty"`
	Raw     *string                `json:"raw,omitempty"`
}

// Postman stores the type of data associated with this request in this field.
type Mode string

const (
	Graphql Mode = "graphql"
	Raw     Mode = "raw"
)
