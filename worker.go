package hashgrab

import (
	"runtime"
	"sync"
)

type Worker struct {
	parallel int
	fetcher  Fetcher
	hasher   Hasher
	sync.WaitGroup
}

type Response struct {
	Url   string
	Hash  string
	Error error
}

func New() *Worker {
	return &Worker{
		parallel: runtime.NumCPU(),
		fetcher:  NewFetcher(),
		hasher:   md5Hasher{},
	}
}

func (w *Worker) MaxWorker(n int) *Worker {
	w.parallel = n
	return w
}

func (w *Worker) Fetcher(fetcher Fetcher) *Worker {
	w.fetcher = fetcher
	return w
}

func (w *Worker) Hasher(hasher Hasher) *Worker {
	w.hasher = hasher
	return w
}

func (w *Worker) Run(urls []string) chan *Response {
	ch := make(chan *Response)
	sem := NewSemaphore(w.parallel)

	go func() {
		defer close(ch)
		for _, url := range urls {
			w.Add(1)
			sem.Acquire()
			go w.process(url, ch, sem)
		}
		w.Wait()
	}()

	return ch
}

func (w *Worker) process(url string, ch chan *Response, sem *Semaphore) {
	defer func() {
		sem.Release()
		w.Done()
	}()

	data, err := w.fetcher.Fetch(url)
	if err != nil {
		ch <- &Response{
			Url:   url,
			Error: err,
		}
		return
	}
	hash := w.hasher.Hash(data)
	ch <- &Response{
		Url:  url,
		Hash: hash,
	}
}
