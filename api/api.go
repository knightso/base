package api

// http://tools.ietf.org/html/draft-nottingham-http-problem-06
type ErrorResponse struct {
	Status     int         `json:"status"`
	Type       string      `json:"type,omitempty"`
	Title      string      `json:"title,omitempty"`
	Detail     string      `json:"detail,omitempty"`
	Instance   string      `json:"instance,omitempty"`
	Additional interface{} `json:"additional,omitempty"`
}

type HasAdditional interface {
	Additional() interface{}
}
