package hashgrab

import "context"

// Semaphore is a simple implementation of a counting semaphore using a channel.
type Semaphore struct {
	// sem is a channel that will be used to manage resources.
	sem chan struct{}
}

// NewSemaphore initializes a new Semaphore with a provided limit.
// It returns a pointer to the Semaphore instance.
func NewSemaphore(limit int) *Semaphore {
	if limit <= 0 {
		panic("hashgrab: semaphore limit must be greater than zero")
	}
	// The channel is buffered with a capacity of 'limit'.
	// Any write to the channel blocks when 'limit' number of writes have been made
	// that have not yet been matched by reads.
	return &Semaphore{sem: make(chan struct{}, limit)}
}

// Acquire acquires a unit of resource.
// If all units of the resource are occupied, Acquire blocks until a unit becomes free or the context is cancelled.
func (s *Semaphore) Acquire(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// Writing to the 'sem' channel represents acquiring a unit of resource.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.sem <- struct{}{}:
		return nil
	}
}

// Release releases a unit of resource.
func (s *Semaphore) Release() {
	// Reading from the 'sem' channel represents releasing a unit of resource.
	<-s.sem
}
