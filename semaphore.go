package hashgrab

type Semaphore struct {
	sem chan struct{}
}

func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{sem: make(chan struct{}, limit)}
}

func (s *Semaphore) Acquire() {
	s.sem <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.sem
}
