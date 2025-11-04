package hashgrab

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Fetcher is an interface that defines methods for fetching data from a URL.
type Fetcher interface {
	// Fetch takes a context and a URL and returns the fetched data as a byte slice, or an error if the fetching failed.
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// Hasher is an interface that defines methods for hashing data.
type Hasher interface {
	// Hash takes a byte slice and returns its hash as a string.
	Hash(data []byte) string
}

// httpFetcher is a struct that implements the Fetcher interface using HTTP.
type httpFetcher struct {
	client *http.Client
}

// NewFetcher returns a new Fetcher that fetches data using HTTP.
func NewFetcher() Fetcher {
	return &httpFetcher{client: &http.Client{Timeout: 15 * time.Second}}
}

// Fetch fetches data from the given URL using HTTP. If the URL does not have a "http://" or "https://" prefix,
// "http://" is added. It returns the fetched data as a byte slice, or an error if the fetching failed.
func (f *httpFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request for %s: %w", url, err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("fetch %s: unexpected status %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", url, err)
	}
	return body, nil
}

// md5Hasher is a struct that implements the Hasher interface using the MD5 algorithm.
type md5Hasher struct{}

// NewMD5Hasher returns a new Hasher that hashes provided data using md5.
func NewMD5Hasher() Hasher {
	return md5Hasher{}
}

// Hash hashes the given data using the MD5 algorithm and returns the hash as a string.
func (md5Hasher) Hash(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}

// sha256Hasher is a struct that implements the Hasher interface using the SHA-256 algorithm.
type sha256Hasher struct{}

// NewSHA256Hasher returns a new Hasher that hashes provided data using sha256.
func NewSHA256Hasher() Hasher {
	return sha256Hasher{}
}

// Hash hashes the given data using the SHA-256 algorithm and returns the hash as a string.
func (sha256Hasher) Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// NewHasherFromName returns a Hasher implementation based on the provided algorithm name.
// Supported algorithms are "md5" and "sha256".
func NewHasherFromName(name string) (Hasher, error) {
	switch strings.ToLower(name) {
	case "md5":
		return NewMD5Hasher(), nil
	case "sha256", "":
		return NewSHA256Hasher(), nil
	default:
		return nil, fmt.Errorf("unsupported hash algorithm %q", name)
	}
}
