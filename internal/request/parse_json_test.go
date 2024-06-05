package request_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/internal/request"
)

func TestParseJSON(t *testing.T) {
	l, tp, m, src, f := 1, "type-text", "message-text", "source-text", "fingerprint-text"
	s := []byte(fmt.Sprintf(`{"level":%d,"type":"%s","message":"%s","source":"%s","fingerprint":"%s"}`, l, tp, m, src, f))

	r, err := request.ParseJSON(s)
	require.NoError(t, err)

	assert.Equal(t, r.Level, request.Level(l))
	assert.Equal(t, r.Type, request.Type(tp))
	assert.Equal(t, r.Message, request.Message(m))
	assert.Equal(t, r.Source, request.Source(src))
	assert.Equal(t, r.Fingerprint, request.Fingerprint(f))
}
