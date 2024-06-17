package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/internal/util"
)

func TestIntToMillisecond(t *testing.T) {
	ms := 500
	duration := 500 * time.Millisecond

	require.Equal(t, duration, util.IntToMillisecond(ms))
}
