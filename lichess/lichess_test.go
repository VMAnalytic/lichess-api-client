package lichess

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	baseURLPath = "/api"
)

func setUp() (client *Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Client.baseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	client = NewClient("API_KEY", nil)
	uri, _ := url.Parse(server.URL + baseURLPath + "/")
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
