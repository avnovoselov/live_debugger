package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/internal/server/internal"
)

func TestBuildAddress(t *testing.T) {
	ip, port, expectedAddress := "0.0.0.0", "80", "0.0.0.0:80"

	address := internal.BuildAddress(ip, port)
	require.Equal(t, expectedAddress, address)
}
