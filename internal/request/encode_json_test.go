package request_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/avnovoselov/live_debugger/internal/request"
)

func TestEncodeJSON(t *testing.T) {
	var expectedID uint64 = 123
	var expectedErr string = "error"

	r1 := request.LogResponse{Offset: &expectedID}
	j1, err := request.EncodeJSON(r1)

	r2 := request.LogResponse{}
	j2, err := request.EncodeJSON(r2)

	r3 := request.LogResponse{Error: &expectedErr}
	j3, err := request.EncodeJSON(r3)

	assert.NoError(t, err)
	assert.Equal(t, `{"offset":123,"error":null}`, string(j1))

	assert.NoError(t, err)
	assert.Equal(t, `{"offset":null,"error":null}`, string(j2))

	assert.NoError(t, err)
	assert.Equal(t, `{"offset":null,"error":"error"}`, string(j3))
}
