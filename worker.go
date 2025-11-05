package hashgrab

import (
	"context"
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
		hasher:   NewSHA256Hasher(),
	}
}

// MaxWorker sets the maximum number of workers to n and returns the updated Worker.
func (w *Worker) MaxWorker(n int) *Worker {
	if n <= 0 {
		panic("hashgrab: max worker must be greater than zero")
	}
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
	return w.RunContext(context.Background(), urls)
}

// RunContext starts fetching and hashing operation on the provided list of urls respecting the given context.
// It returns a channel of Response where the results of the operations are sent.
func (w *Worker) RunContext(ctx context.Context, urls []string) chan *Response {
	// Channel to collect the results.
	ch := make(chan *Response)
	// Semaphore to limit the number of concurrent operations.
	sem := NewSemaphore(w.parallel)
	var wg sync.WaitGroup

	// Start a goroutine to manage the operations.
	go func() {
		// Close the results channel when all operations are done.
		defer close(ch)
		var started bool
		// Loop over the URLs.
		for _, url := range urls {
			if err := ctx.Err(); err != nil {
				break
			}
			if err := sem.Acquire(ctx); err != nil {
				break
			}
			if err := ctx.Err(); err != nil {
				sem.Release()
				break
			}
			// For each URL, add to the wait group and start the process in a separate goroutine.
			wg.Add(1)
			started = true
			go w.process(ctx, url, ch, sem, &wg)
		}
		if !started {
			return
		}
		// Wait for all operations to complete.
		wg.Wait()
	}()

	// Return the results channel.
	return ch
}

// process is a helper function that fetches and hashes data from a URL,
// sends the result on a channel and releases a semaphore.
func (w *Worker) process(ctx context.Context, url string, ch chan *Response, sem *Semaphore, wg *sync.WaitGroup) {
	defer func() {
		// Release the semaphore and signal completion to the wait group when done.
		sem.Release()
		wg.Done()
	}()

	// Fetch data from the URL.
	data, err := w.fetcher.Fetch(ctx, url)
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
