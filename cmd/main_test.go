package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
)

func TestRun_NoURLs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(context.Background(), &stdout, &stderr, []string{"-parallel", "1"})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Please provide at least one URL") {
		t.Fatalf("expected error message about URLs, got %q", stderr.String())
	}
}

func TestRun_InvalidParallel(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(context.Background(), &stdout, &stderr, []string{"-parallel", "0", "http://example.com"})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Number of parallel requests should be greater than 0") {
		t.Fatalf("expected error about parallel value, got %q", stderr.String())
	}
}

func TestRun_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("payload"))
	}))
	defer server.Close()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), &stdout, &stderr, []string{"-parallel", "2", server.URL})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr=%q)", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), server.URL) {
		t.Fatalf("expected output to contain URL, got %q", stdout.String())
	}
	if stdout.Len() == 0 {
		t.Fatalf("expected hash output, got empty string")
	}
}
