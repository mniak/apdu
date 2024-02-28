package test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bytesOrHex interface{ ~[]byte | string }

func bytesFrom[T bytesOrHex](value T) ([]byte, error) {
	if str, ok := (any(value)).(string); ok {
		return hex.DecodeString(str)
	}

	if byt, ok := (any(value)).([]byte); ok {
		return byt, nil
	}
	panic("invalid type")
}

func AssertBytesEqual[E, A bytesOrHex](t *testing.T, expected E, actual A, msgAndArgs ...any) {
	t.Helper()

	expectedBytes, err := bytesFrom(expected)
	require.NoError(t, err)

	actualBytes, err := bytesFrom(actual)
	require.NoError(t, err)

	assert.Equal(t,
		fmt.Sprintf("% 02X", expectedBytes),
		fmt.Sprintf("% 02X", actualBytes),
		msgAndArgs,
	)
}
