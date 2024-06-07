package util_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/internal/util"
)

func TestErrorBuilder(t *testing.T) {
	baseError, extendedError := errors.New("base"), errors.New("extended")

	eb := util.ErrorBuilder(baseError)
	resultError := eb(extendedError)

	require.ErrorIs(t, resultError, baseError)
	require.ErrorIs(t, resultError, extendedError)
}
