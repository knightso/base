package api

// http://tools.ietf.org/html/draft-nottingham-http-problem-06
type ErrorResponse struct {
	Status   int    `json:"status"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}
