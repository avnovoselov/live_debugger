package queue

import (
	"errors"
	"sync"

	"github.com/avnovoselov/live_debugger/internal/util"
)

// queueEB - queue error builder contains base error
var queueEB = util.ErrorBuilder(errors.New("queue error"))

var (
	OffsetNotFoundError = queueEB(errors.New("offset not found"))
	QueueIsEmpty        = queueEB(errors.New("queue is empty"))
)

// Queue - sized queue.
// All public operations:
//
// - GetAll
// - GetByOffset
// - Append
//
// execs synchronously using sync.RWMutex to prevent race condition.
type Queue[Element any] struct {
	// queue dimensions helpers
	min  uint64
	max  uint64
	size uint64

	// queued elements
	elements map[uint64]Element

	// r/w lock
	mu sync.RWMutex
}

// NewQueue - Queue constructor.
func NewQueue[Element any](size uint64) *Queue[Element] {
	return &Queue[Element]{
		size:     size,
		elements: make(map[uint64]Element),
		mu:       sync.RWMutex{},
	}
}

// GetAll - return all queued elements.
func (q *Queue[Element]) GetAll() []Element {
	q.mu.RLock()
	defer q.mu.RUnlock()

	r := make([]Element, 0, len(q.elements))

	for i := q.min; i < q.max; i++ {
		r = append(r, q.elements[i])
	}

	return r
}

// GetByOffset - return an element by its offset.
// If the offset too old and not represent at the queue returns the oldest element and its offset.
// If the offset gte q.max returns error and last element offset.
func (q *Queue[Element]) GetByOffset(offset uint64) (Element, uint64, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if offset < q.min {
		return q.elements[q.min], q.min, nil
	}
	if offset >= q.max {
		return *new(Element), q.max - 1, OffsetNotFoundError
	}

	return q.elements[offset], offset, nil
}

// GetLast - returns the last element or error if queue is empty
func (q *Queue[Element]) GetLast() (Element, uint64, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.elements) == 0 {
		return *new(Element), 0, QueueIsEmpty
	}

	return q.elements[q.max-1], q.max - 1, nil
}

// Append - put an element to the queues tail and delete oldest
func (q *Queue[Element]) Append(element Element) uint64 {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.appendElement(element)
	q.clearOldest()

	// Should use minus one to define the last existing element not a next element
	return q.max - 1
}

// appendElement - places an element to the maps q.max position and increment q.max.
func (q *Queue[Element]) appendElement(element Element) {
	q.elements[q.max] = element
	q.max += 1
}

// clearOldest - deletes a maps q.min element if the queue size exceeded.
func (q *Queue[Element]) clearOldest() {
	if len(q.elements) > int(q.size) {
		delete(q.elements, q.min)

		q.min += 1
	}
}
