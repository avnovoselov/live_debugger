package queue_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/internal/queue"
)

func TestQueue_Append_Size(t *testing.T) {
	var (
		size  uint64 = 30
		retry        = 1000
	)

	q := queue.NewQueue[int](size)
	wg := sync.WaitGroup{}

	for i := 0; i < retry; i++ {
		wg.Add(1)
		go func(i int) {
			q.Append(i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	require.Len(t, q.GetAll(), int(size))
}

func TestQueue_GetAll(t *testing.T) {
	var size uint64 = 3

	q := queue.NewQueue[int](size)

	q.Append(1)
	assert.Equal(t, []int{1}, q.GetAll())
	q.Append(2)
	assert.Equal(t, []int{1, 2}, q.GetAll())
	q.Append(3)
	assert.Equal(t, []int{1, 2, 3}, q.GetAll())
	q.Append(4)
	assert.Equal(t, []int{2, 3, 4}, q.GetAll())
	q.Append(5)
	assert.Equal(t, []int{3, 4, 5}, q.GetAll())
}

func TestQueue_GetByOffset(t *testing.T) {
	type TestCase struct {
		name              string
		offset            uint64
		expectedNewOffset uint64
		expectedElement   int
		expectedErr       error
	}

	testCases := []TestCase{
		{
			name:              "First element",
			offset:            2,
			expectedNewOffset: 3,
			expectedElement:   3,
			expectedErr:       nil,
		},
		{
			name:              "Last element",
			offset:            5,
			expectedNewOffset: 5,
			expectedElement:   5,
			expectedErr:       nil,
		},
		{
			name:              "Old missed element",
			offset:            0,
			expectedNewOffset: 3,
			expectedElement:   3,
			expectedErr:       nil,
		},
		{
			name:              "Not exists offset",
			offset:            100,
			expectedNewOffset: 5,
			expectedElement:   0,
			expectedErr:       queue.OffsetNotFoundError,
		},
	}

	q := queue.NewQueue[int](3)
	q.Append(1)
	q.Append(2)
	q.Append(3)
	q.Append(4)
	q.Append(5)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			element, newOffset, err := q.GetByOffset(testCase.offset)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Equal(t, testCase.expectedNewOffset, newOffset)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedNewOffset, newOffset)
				assert.Equal(t, testCase.expectedElement, element)
			}
		})
	}
}

func TestQueue_GetLast(t *testing.T) {
	q := queue.NewQueue[int](30)
	_, _, err := q.GetLast()
	require.ErrorIs(t, err, queue.IsEmpty)

	first, second, third := 1, 2, 3
	q.Append(first)

	e, o, err := q.GetLast()
	require.NoError(t, err)
	assert.Equal(t, first, e)
	assert.Equal(t, uint64(1), o)

	q.Append(second)
	e, o, err = q.GetLast()
	require.NoError(t, err)
	assert.Equal(t, second, e)
	assert.Equal(t, uint64(2), o)

	q.Append(third)
	e, o, err = q.GetLast()
	require.NoError(t, err)
	assert.Equal(t, third, e)
	assert.Equal(t, uint64(3), o)

	e, o, err = q.GetLast()
	require.NoError(t, err)
	assert.Equal(t, third, e)
	assert.Equal(t, uint64(3), o)
}
