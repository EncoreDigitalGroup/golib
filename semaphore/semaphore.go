/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Rights Reserved.
 */

package semaphore

import (
	"errors"
	"runtime"
	"time"
)

var ErrAcquireTimeout = errors.New("failed to acquire semaphore: context timeout or cancellation")
var ErrNoSlotsAvailable = errors.New("failed to acquire semaphore: no slots available")

type Semaphore struct {
	Workers int
	Channel chan struct{}
}

func Make(workers *int) *Semaphore {
	workerCount := runtime.NumCPU()

	if workers != nil {
		workerCount = *workers
	}

	return &Semaphore{
		Workers: workerCount,
		Channel: make(chan struct{}, workerCount),
	}
}

// Lock acquires a semaphore slot, executes the function, and releases the slot.
// If timeout is nil, this method blocks indefinitely until a slot is available.
// Otherwise, it returns ErrAcquireTimeout if the timeout expires before acquiring a slot.
func (s *Semaphore) Lock(fn func() error, timeout *time.Duration) error {
	if timeout == nil {
		s.Channel <- struct{}{}
		defer func() { <-s.Channel }()
		return fn()
	}

	select {
	case s.Channel <- struct{}{}:
		defer func() { <-s.Channel }()
		return fn()
	case <-time.After(*timeout):
		return ErrAcquireTimeout
	}
}

// TryLock attempts to acquire a semaphore slot without blocking.
// If a slot is available, it executes the function and returns its result.
// If no slot is available, it returns ErrNoSlotsAvailable immediately without executing the function.
func (s *Semaphore) TryLock(fn func() error) error {
	select {
	case s.Channel <- struct{}{}:
		defer func() { <-s.Channel }()
		return fn()
	default:
		return ErrNoSlotsAvailable
	}
}

// AvailableSlots returns the number of currently available slots in the semaphore
func (s *Semaphore) AvailableSlots() int {
	return cap(s.Channel) - len(s.Channel)
}

// TotalSlots returns the total number of slots in the semaphore
func (s *Semaphore) TotalSlots() int {
	return s.Workers
}

// WaitForCompletion blocks until all semaphore operations complete.
func (s *Semaphore) WaitForCompletion() error {
	tokens := make([]struct{}, 0, s.Workers)

	releaseTokens := func() {
		for range tokens {
			<-s.Channel
		}
	}

	for i := 0; i < s.Workers; i++ {
		s.Channel <- struct{}{}
		tokens = append(tokens, struct{}{})
	}
	releaseTokens()
	return nil
}
