// Package hashgrab (in this file) provides functionality for fetching data from a URL and hashing it.
// It is used in worker.go to fetch data from URLs and hash it.
package hashgrab

import (
	"fmt"
	"io"

	"crypto/md5"
	"net/http"
	"strings"
)

// Fetcher is an interface that defines methods for fetching data from a URL.
type Fetcher interface {
	// Fetch takes a URL and returns the fetched data as a byte slice, or an error if the fetching failed.
	Fetch(url string) ([]byte, error)
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
	return &httpFetcher{client: &http.Client{}}
}

// Fetch fetches data from the given URL using HTTP. If the URL does not have a "http://" or "https://" prefix,
// "http://" is added. It returns the fetched data as a byte slice, or an error if the fetching failed.
func (f *httpFetcher) Fetch(url string) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	resp, err := f.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
