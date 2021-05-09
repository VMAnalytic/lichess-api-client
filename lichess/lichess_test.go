package lichess

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	baseURLPath = "/"
)

func setUp() (client *Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath, mux)

	server := httptest.NewServer(apiHandler)

	client = NewClient("API_KEY", nil)
	uri, _ := url.Parse(server.URL + "/")
	client.baseURL = uri

	return client, mux, server.Close
}

//nolint
func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()

	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}
