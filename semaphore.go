// Package hashgrab (in this file) provides a basic implementation of a counting semaphore
// for managing concurrency. It is used in worker.go to limit the number of concurrent.
package hashgrab

// Semaphore is a simple implementation of a counting semaphore using a channel.
type Semaphore struct {
	// sem is a channel that will be used to manage resources.
	sem chan struct{}
}

// NewSemaphore initializes a new Semaphore with a provided limit.
// It returns a pointer to the Semaphore instance.
func NewSemaphore(limit int) *Semaphore {
	// The channel is buffered with a capacity of 'limit'.
	// Any write to the channel blocks when 'limit' number of writes have been made
	// that have not yet been matched by reads.
	return &Semaphore{sem: make(chan struct{}, limit)}
}

// Acquire acquires a unit of resource.
// If all units of the resource are occupied, Acquire blocks until a unit becomes free.
func (s *Semaphore) Acquire() {
	// Writing to the 'sem' channel represents acquiring a unit of resource.
	s.sem <- struct{}{}
}

// Release releases a unit of resource.
func (s *Semaphore) Release() {
	// Reading from the 'sem' channel represents releasing a unit of resource.
	<-s.sem
}
