package live_debugger_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/pkg/live_debugger"
)

func TestEncodeJSON(t *testing.T) {
	var expectedID uint64 = 123
	var expectedErr = "error"

	r1 := live_debugger.LogCreatedDTO{Offset: &expectedID}
	j1, err := live_debugger.EncodeJSON(r1)
	require.NoError(t, err)

	r2 := live_debugger.LogCreatedDTO{}
	j2, err := live_debugger.EncodeJSON(r2)
	require.NoError(t, err)

	r3 := live_debugger.LogCreatedDTO{Error: &expectedErr}
	j3, err := live_debugger.EncodeJSON(r3)
	require.NoError(t, err)

	assert.Equal(t, `{"offset":123,"error":null}`, string(j1))
	assert.Equal(t, `{"offset":null,"error":null}`, string(j2))
	assert.Equal(t, `{"offset":null,"error":"error"}`, string(j3))
}

func TestParseJSON(t *testing.T) {
	l, tp, m, src, f := 1, "type-text", "message-text", "source-text", "fingerprint-text"
	s := []byte(fmt.Sprintf(`{"level":%d,"type":"%s","message":"%s","source":"%s","fingerprint":"%s"}`, l, tp, m, src, f))

	r, err := live_debugger.ParseJSON[live_debugger.LogDTO](s)
	require.NoError(t, err)

	assert.Equal(t, r.Level, live_debugger.Level(l))
	assert.Equal(t, r.Type, live_debugger.Type(tp))
	assert.Equal(t, r.Message, live_debugger.Message(m))
	assert.Equal(t, r.Source, live_debugger.Source(src))
	assert.Equal(t, r.Fingerprint, live_debugger.Fingerprint(f))
}
