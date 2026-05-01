package hashgrab

import (
	"context"
	"errors"
	"runtime"
	"testing"
)

type mockFetcher struct {
	FetchFunc func(context.Context, string) ([]byte, error)
}

func (m mockFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	if m.FetchFunc != nil {
		return m.FetchFunc(ctx, url)
	}
	return []byte{}, nil
}

func getMockFetcher() Fetcher {
	return mockFetcher{
		FetchFunc: func(ctx context.Context, url string) ([]byte, error) {
			return []byte("testdata"), nil
		},
	}
}

type mockHasher struct {
	HashFunc func([]byte) string
}

func (m mockHasher) Hash(data []byte) string {
	if m.HashFunc != nil {
		return m.HashFunc(data)
	}
	return ""
}

func getMockHasher() Hasher {
	return mockHasher{
		HashFunc: func(data []byte) string {
			return "testHash"
		},
	}
}

func TestWorker_New(t *testing.T) {
	worker := New()
	if worker.parallel != runtime.NumCPU() {
		t.Errorf("expected parallel to be %d, got %d", runtime.NumCPU(), worker.parallel)
	}
	if _, ok := worker.fetcher.(*httpFetcher); !ok {
		t.Errorf("expected fetcher to be of type *httpFetcher, got %T", worker.fetcher)
	}
	if _, ok := worker.hasher.(sha256Hasher); !ok {
		t.Errorf("expected hasher to be of type sha256Hasher, got %T", worker.hasher)
	}
}

func TestWorker_MaxWorker(t *testing.T) {
	worker := New()
	worker.MaxWorker(5)
	if worker.parallel != 5 {
		t.Errorf("expected parallel to be %d, got %d", 5, worker.parallel)
	}
}

func TestWorker_MaxWorkerInvalid(t *testing.T) {
	worker := New()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when setting non-positive worker count")
		}
	}()
	worker.MaxWorker(0)
}

func TestWorker_Fetcher(t *testing.T) {
	worker := New()
	testFetcher := &mockFetcher{}
	worker.Fetcher(testFetcher)
	if worker.fetcher != testFetcher {
		t.Errorf("expected fetcher to be %p, got %p", testFetcher, worker.fetcher)
	}
}

func TestWorker_Hasher(t *testing.T) {
	worker := New()
	testHasher := &mockHasher{}
	worker.Hasher(testHasher)
	if worker.hasher != testHasher {
		t.Errorf("expected hasher to be %p, got %p", testHasher, worker.hasher)
	}
}

func TestWorker_Run(t *testing.T) {
	fetcher, hasher := getMockFetcher(), getMockHasher()
	tests := []struct {
		name        string
		urls        []string
		fetcher     Fetcher
		hasher      Hasher
		maxWorkers  int
		expectedOut []*Response
	}{
		{
			name:       "Single URL",
			urls:       []string{"http://example.com"},
			fetcher:    fetcher,
			hasher:     hasher,
			maxWorkers: 10,
			expectedOut: []*Response{
				{
					Url:  "http://example.com",
					Hash: "testHash",
				},
			},
		},
		{
			name:       "Multiple URLs",
			urls:       []string{"http://example.com", "http://test.com"},
			fetcher:    fetcher,
			hasher:     hasher,
			maxWorkers: 1,
			expectedOut: []*Response{
				{
					Url:  "http://example.com",
					Hash: "testHash",
				},
				{
					Url:  "http://test.com",
					Hash: "testHash",
				},
			},
		},
		{
			name: "Fetch Error",
			urls: []string{"http://example.com"},
			fetcher: mockFetcher{
				FetchFunc: func(ctx context.Context, url string) ([]byte, error) {
					return nil, errors.New("fetch error")
				},
			},
			hasher:     hasher,
			maxWorkers: 10,
			expectedOut: []*Response{
				{
					Url:   "http://example.com",
					Error: errors.New("fetch error"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := New().MaxWorker(tt.maxWorkers).Fetcher(tt.fetcher).Hasher(tt.hasher)
			outCh := worker.Run(tt.urls)
			for _, expected := range tt.expectedOut {
				out := <-outCh
				if out.Url != expected.Url {
					t.Errorf("expected URL %q, got %q", expected.Url, out.Url)
				}
				if out.Hash != expected.Hash {
					t.Errorf("expected hash %q, got %q", expected.Hash, out.Hash)
				}
				if (out.Error != nil) != (expected.Error != nil) {
					t.Errorf("expected error %v, got %v", expected.Error, out.Error)
				}
			}
		})
	}
}

func TestWorker_RunContextCanceled(t *testing.T) {
	worker := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	out := worker.Fetcher(getMockFetcher()).Hasher(getMockHasher()).RunContext(ctx, []string{"http://example.com"})

	resp, ok := <-out
	if !ok {
		t.Fatalf("expected cancellation response")
	}
	if resp.Url != "http://example.com" {
		t.Fatalf("expected URL %q, got %q", "http://example.com", resp.Url)
	}
	if resp.Error == nil {
		t.Fatalf("expected error due to cancellation")
	}
	if !errors.Is(resp.Error, context.Canceled) {
		t.Fatalf("expected context cancellation error, got %v", resp.Error)
	}
	if _, ok := <-out; ok {
		t.Fatalf("expected channel to close after single response")
	}
}
