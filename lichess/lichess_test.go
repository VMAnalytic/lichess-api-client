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

func TestClient_SetLimits(t *testing.T) {
	client := NewClient("API_KEY", nil)

	err := client.SetLimits(1, 5)

	if err != nil {
		t.Errorf("Set limits should not throw an error: %v", err)
	}
}

func TestNewClient_withDefaultHTTPClient(t *testing.T) {
	c := NewClient("", nil)

	if got, want := c.baseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}

	c2 := NewClient("nil", nil)
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}
