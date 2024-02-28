package test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanHex(hexString string) string {
	hexString = strings.ReplaceAll(hexString, " ", "")
	hexString = strings.ReplaceAll(hexString, "_", "")
	hexString = strings.ReplaceAll(hexString, "-", "")
	hexString = strings.ReplaceAll(hexString, "\n", "")
	hexString = strings.ReplaceAll(hexString, "\t", "")
	return hexString
}

func ParseHex(t *testing.T, hexString string) []byte {
	t.Helper()
	hexString = cleanHex(hexString)
	data, err := hex.DecodeString(hexString)
	assert.NoError(t, err)
	return data
}

func MustParseHex(t *testing.T, hexString string) []byte {
	t.Helper()
	hexString = cleanHex(hexString)
	data, err := hex.DecodeString(hexString)
	require.NoError(t, err)
	return data
}
