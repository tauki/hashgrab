package hashgrab

import (
	"runtime"
	"sync"
)

// Worker struct orchestrates the fetching and hashing of URL data.
type Worker struct {
	// parallel is the maximum number of concurrent fetch and hash operations.
	parallel int
	// fetcher is the interface for fetching data from a URL.
	fetcher Fetcher
	// hasher is the interface for hashing the fetched data.
	hasher Hasher
	// WaitGroup is embedded to track completion of all fetch and hash operations.
	sync.WaitGroup
}

// Response struct represents the result of a fetch and hash operation.
type Response struct {
	Url   string
	Hash  string
	Error error
}

// New creates a new Worker with default configuration.
// It initializes fetcher and hasher with their default implementations
// and sets the number of parallel workers to the number of CPUs.
func New() *Worker {
	return &Worker{
		parallel: runtime.NumCPU(),
		fetcher:  NewFetcher(),
		hasher:   NewMD5Hasher(),
	}
}

// MaxWorker sets the maximum number of workers to n and returns the updated Worker.
func (w *Worker) MaxWorker(n int) *Worker {
	w.parallel = n
	return w
}

// GetMaxWorker returns the maximum number of workers.
func (w *Worker) GetMaxWorker() int {
	return w.parallel
}

// Fetcher sets the fetcher implementation to fetcher and returns the updated Worker.
func (w *Worker) Fetcher(fetcher Fetcher) *Worker {
	w.fetcher = fetcher
	return w
}

// Hasher sets the hasher implementation to hasher and returns the updated Worker.
func (w *Worker) Hasher(hasher Hasher) *Worker {
	w.hasher = hasher
	return w
}

// Run starts fetching and hashing operation on the provided list of urls.
// It returns a channel of Response where the results of the operations are sent.
func (w *Worker) Run(urls []string) chan *Response {
	// Channel to collect the results.
	ch := make(chan *Response)
	// Semaphore to limit the number of concurrent operations.
	sem := NewSemaphore(w.parallel)

	// Start a goroutine to manage the operations.
	go func() {
		// Close the results channel when all operations are done.
		defer close(ch)
		// Loop over the URLs.
		for _, url := range urls {
			// For each URL, add to the wait group and acquire a semaphore.
			w.Add(1)
			sem.Acquire()
			// Start the process in a separate goroutine.
			go w.process(url, ch, sem)
		}
		// Wait for all operations to complete.
		w.Wait()
	}()

	// Return the results channel.
	return ch
}

// process is a helper function that fetches and hashes data from a URL,
// sends the result on a channel and releases a semaphore.
func (w *Worker) process(url string, ch chan *Response, sem *Semaphore) {
	defer func() {
		// Release the semaphore and signal completion to the wait group when done.
		sem.Release()
		w.Done()
	}()

	// Fetch data from the URL.
	data, err := w.fetcher.Fetch(url)
	if err != nil {
		// Send an error response and return if fetching failed.
		ch <- &Response{
			Url:   url,
			Error: err,
		}
		return
	}
	// Hash the fetched data.
	hash := w.hasher.Hash(data)
	// Send the successful response.
	ch <- &Response{
		Url:  url,
		Hash: hash,
	}
}
