package http

// HTTPResponse represents a standard HTTP response
type HTTPResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}