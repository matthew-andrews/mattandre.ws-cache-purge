package cloudflare

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPurge(t *testing.T) {

	// Make a server that always responds with a successfully purged response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{ "result": { "id": "123" }, "success": true }`)
	}))
	defer server.Close()

	// Make a transport that reroutes all traffic to the example server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	// Make a HTTP client with that transport
	httpClient := &http.Client{Transport: transport}

	// Make an API client and *inject* (it shouldn't hurt…)
	client := &Client{server.URL, httpClient}

	// Test the method!
	resp, err := client.Purge([]string{"https://mattandre.ws/index.html", "https://mattandre.ws/"})

	if err != nil {
		t.Fatalf("expected Purge not to error but it did with %v", err)
	}

	if resp.Success != true {
		t.Fatalf("expect success to be true, got %v", resp.Success)
	}
}
