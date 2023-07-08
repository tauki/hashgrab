package hashgrab_test

import (
	"fmt"

	"crypto/sha256"
	"github.com/tauki/hashgrab"
)

type customFetcher struct{}

func (f *customFetcher) Fetch(url string) ([]byte, error) {
	// Implement custom fetching logic here
	// This is a dummy implementation that returns URL as bytes and nil error
	return []byte(url), nil
}

type customHasher struct{}

func (h *customHasher) Hash(data []byte) string {
	// Implement custom hashing logic here
	// This is a dummy implementation that returns SHA256 hash as hex
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// ExampleNew demonstrates how to create a new Worker.
func ExampleNew() {
	worker := hashgrab.New()
	fmt.Printf("Created a new worker with %v parallelism\n", worker.GetMaxWorker())
	// Output: Created a new worker with 10 parallelism
}

// ExampleWorker_MaxWorker demonstrates how to set the maximum number of workers.
func ExampleWorker_MaxWorker() {
	worker := hashgrab.New().MaxWorker(5)
	fmt.Printf("Updated worker to have %v parallelism\n", worker.GetMaxWorker())
	// Output: Updated worker to have 5 parallelism
}

// ExampleWorker_Fetcher demonstrates how to set a custom fetcher.
func ExampleWorker_Fetcher() {
	worker := hashgrab.New()
	customFetcher := hashgrab.NewFetcher() // Assumes you have a custom Fetcher
	worker.Fetcher(customFetcher)
	fmt.Println("Updated worker to use a custom fetcher")
	// Output: Updated worker to use a custom fetcher
}

// ExampleWorker_Hasher demonstrates how to set a custom hasher.
func ExampleWorker_Hasher() {
	worker := hashgrab.New()
	customHasher := hashgrab.NewMD5Hasher() // Assumes you have a custom Hasher
	worker.Hasher(customHasher)
	fmt.Println("Updated worker to use a custom hasher")
	// Output: Updated worker to use a custom hasher
}

// ExampleWorker_Run demonstrates how to start a Worker.
func ExampleWorker_Run() {
	worker := hashgrab.New().MaxWorker(2)
	worker.Fetcher(&customFetcher{})
	worker.Hasher(&customHasher{})
	urls := []string{"http://example.com", "http://example.org"}
	responseChan := worker.Run(urls)

	for response := range responseChan {
		if response.Error != nil {
			fmt.Printf("Error fetching %s: %v\n", response.Url, response.Error)
			continue
		}
		fmt.Printf("Hash of %s: %s\n", response.Url, response.Hash)
	}
	// Output:
	// Hash of http://example.org: 971a565c8ac770ff0b288d98a507bb832b8002214411ed8244d0b981a506dd3e
	// Hash of http://example.com: f0e6a6a97042a4f1f1c87f5f7d44315b2d852c2df5c7991cc66241bf7072d1c4
}
