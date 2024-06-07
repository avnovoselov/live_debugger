package response_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/internal/response"
)

func TestEncodeJSON(t *testing.T) {
	var expectedID uint64 = 123
	var expectedErr string = "error"

	r1 := response.LogResponse{Offset: &expectedID}
	j1, err := response.EncodeJSON(r1)
	require.NoError(t, err)

	r2 := response.LogResponse{}
	j2, err := response.EncodeJSON(r2)
	require.NoError(t, err)

	r3 := response.LogResponse{Error: &expectedErr}
	j3, err := response.EncodeJSON(r3)
	require.NoError(t, err)

	assert.Equal(t, `{"offset":123,"error":null}`, string(j1))
	assert.Equal(t, `{"offset":null,"error":null}`, string(j2))
	assert.Equal(t, `{"offset":null,"error":"error"}`, string(j3))
}
