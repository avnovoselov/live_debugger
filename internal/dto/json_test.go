package dto_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/internal/dto"
)

func TestEncodeJSON(t *testing.T) {
	var expectedID uint64 = 123
	var expectedErr = "error"

	r1 := dto.LogCreatedDTO{Offset: &expectedID}
	j1, err := dto.EncodeJSON(r1)
	require.NoError(t, err)

	r2 := dto.LogCreatedDTO{}
	j2, err := dto.EncodeJSON(r2)
	require.NoError(t, err)

	r3 := dto.LogCreatedDTO{Error: &expectedErr}
	j3, err := dto.EncodeJSON(r3)
	require.NoError(t, err)

	assert.Equal(t, `{"offset":123,"error":null}`, string(j1))
	assert.Equal(t, `{"offset":null,"error":null}`, string(j2))
	assert.Equal(t, `{"offset":null,"error":"error"}`, string(j3))
}

func TestParseJSON(t *testing.T) {
	l, tp, m, src, f := 1, "type-text", "message-text", "source-text", "fingerprint-text"
	s := []byte(fmt.Sprintf(`{"level":%d,"type":"%s","message":"%s","source":"%s","fingerprint":"%s"}`, l, tp, m, src, f))

	r, err := dto.ParseJSON[dto.LogDTO](s)
	require.NoError(t, err)

	assert.Equal(t, r.Level, dto.Level(l))
	assert.Equal(t, r.Type, dto.Type(tp))
	assert.Equal(t, r.Message, dto.Message(m))
	assert.Equal(t, r.Source, dto.Source(src))
	assert.Equal(t, r.Fingerprint, dto.Fingerprint(f))
}
