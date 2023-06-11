package hashgrab

import (
	"fmt"
	"io"

	"crypto/md5"
	"net/http"
	"strings"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

type Hasher interface {
	Hash(data []byte) string
}

type httpFetcher struct {
	client *http.Client
}

func NewFetcher() Fetcher {
	return &httpFetcher{client: &http.Client{}}
}

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

type md5Hasher struct{}

func (md5Hasher) Hash(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}
