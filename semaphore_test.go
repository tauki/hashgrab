package hashgrab

import "testing"

func TestNewSemaphore(t *testing.T) {
	sem := NewSemaphore(5)
	if cap(sem.sem) != 5 {
		t.Errorf("NewSemaphore did not set cap correctly, got: %d, want: %d.", cap(sem.sem), 5)
	}
}

func TestAcquireRelease(t *testing.T) {
	sem := NewSemaphore(1)
	sem.Acquire()
	if len(sem.sem) != 1 {
		t.Errorf("Acquire did not increase len correctly, got: %d, want: %d.", len(sem.sem), 1)
	}
	sem.Release()
	if len(sem.sem) != 0 {
		t.Errorf("Release did not decrease len correctly, got: %d, want: %d.", len(sem.sem), 0)
	}
}
