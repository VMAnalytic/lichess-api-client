package lichess

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Response *http.Response
	Message  string         `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

type RateLimitError struct {
	Rate     uint8
	Response *http.Response
	Message  string         `json:"message"`
}

func (r *RateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Rate)
}
