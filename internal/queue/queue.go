package queue

import (
	"errors"
	"sync"
)

var (
	BaseError     = errors.New("queue error")
	NotFoundError = errors.Join(BaseError, errors.New("not found"))
)

// Queue - sized queue.
// All public operations:
//
// - GetAll
// - GetByOffset
// - Append
//
// execs synchronously using sync.RWMutex to prevent race condition.
type Queue[T any] struct {
	// queue dimensions helpers
	min  uint64
	max  uint64
	size uint64

	// queued elements
	elements map[uint64]T

	// r/w lock
	mu sync.RWMutex
}

// NewQueue - Queue constructor.
func NewQueue[T any](size uint64) *Queue[T] {
	return &Queue[T]{
		size:     size,
		elements: make(map[uint64]T),
		mu:       sync.RWMutex{},
	}
}

// GetAll - return all queued elements.
func (q *Queue[T]) GetAll() []T {
	q.mu.RLock()
	defer q.mu.RUnlock()

	r := make([]T, 0, len(q.elements))

	for i := q.min; i < q.max; i++ {
		r = append(r, q.elements[i])
	}

	return r
}

// GetByOffset - return an element by its offset.
// If the offset too old and not represent at the queue returns the oldest element and its offset.
// If the offset gte q.max returns error and last element offset.
func (q *Queue[T]) GetByOffset(offset uint64) (T, uint64, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if offset < q.min {
		return q.elements[q.min], q.min, nil
	}
	if offset >= q.max {
		return *new(T), q.max - 1, NotFoundError
	}

	return q.elements[offset], offset, nil
}

// Append - put an element to the queues tail and delete oldest
func (q *Queue[T]) Append(element T) uint64 {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.appendElement(element)
	q.clearOldest()

	// Should use minus one to define the last existing element not a next element
	return q.max - 1
}

// appendElement - places an element to the maps q.max position and increment q.max.
func (q *Queue[T]) appendElement(element T) {
	q.elements[q.max] = element
	q.max += 1
}

// clearOldest - deletes a maps q.min element if the queue size exceeded.
func (q *Queue[T]) clearOldest() {
	if len(q.elements) > int(q.size) {
		delete(q.elements, q.min)

		q.min += 1
	}
}
