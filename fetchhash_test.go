package hashgrab

import (
	"context"
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
	data, err := fetcher.Fetch(context.Background(), server.URL)
	if err != nil {
		t.Errorf("Fetch returned an error: %s", err)
	}
	if string(data) != "Mock Server Response\n" {
		t.Errorf("Fetch returned incorrect data, got: %s, want: %s.", string(data), "Mock Server Response\n")
	}
}

func TestHash(t *testing.T) {
	hasher := sha256Hasher{}
	hash := hasher.Hash([]byte("test"))
	expectedHash := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if hash != expectedHash {
		t.Errorf("Hash returned incorrect value, got: %s, want: %s.", hash, expectedHash)
	}
}
