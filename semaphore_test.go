package hashgrab

import (
	"context"
	"testing"
	"time"
)

func TestNewSemaphore(t *testing.T) {
	sem := NewSemaphore(5)
	if cap(sem.sem) != 5 {
		t.Errorf("NewSemaphore did not set cap correctly, got: %d, want: %d.", cap(sem.sem), 5)
	}
}

func TestAcquireRelease(t *testing.T) {
	sem := NewSemaphore(1)
	if err := sem.Acquire(context.Background()); err != nil {
		t.Fatalf("Acquire returned unexpected error: %v", err)
	}
	if len(sem.sem) != 1 {
		t.Errorf("Acquire did not increase len correctly, got: %d, want: %d.", len(sem.sem), 1)
	}
	sem.Release()
	if len(sem.sem) != 0 {
		t.Errorf("Release did not decrease len correctly, got: %d, want: %d.", len(sem.sem), 0)
	}
}

func TestAcquireCanceledContext(t *testing.T) {
	sem := NewSemaphore(1)
	// Fill semaphore so acquire blocks.
	if err := sem.Acquire(context.Background()); err != nil {
		t.Fatalf("failed to acquire initial slot: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if err := sem.Acquire(ctx); err == nil {
		t.Fatalf("expected acquire to fail due to context cancellation")
	}
}
