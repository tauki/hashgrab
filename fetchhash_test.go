package hashgrab

import (
	"fmt"
	"testing"

	"net/http"
	"net/http/httptest"
)

func TestFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(rw, "Mock Server Response")
	}))
	defer server.Close()

	fetcher := httpFetcher{client: server.Client()}
	data, err := fetcher.Fetch(server.URL)
	if err != nil {
		t.Errorf("Fetch returned an error: %s", err)
	}
	if string(data) != "Mock Server Response\n" {
		t.Errorf("Fetch returned incorrect data, got: %s, want: %s.", string(data), "Mock Server Response\n")
	}
}

func TestHash(t *testing.T) {
	hasher := md5Hasher{}
	hash := hasher.Hash([]byte("test"))
	expectedHash := "098f6bcd4621d373cade4e832627b4f6"
	if hash != expectedHash {
		t.Errorf("Hash returned incorrect value, got: %s, want: %s.", hash, expectedHash)
	}
}
