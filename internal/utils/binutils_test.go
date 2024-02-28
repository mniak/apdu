package utils

import (
	"testing"

	"github.com/mniak/apdu/internal/test"
	"github.com/stretchr/testify/require"
)

func TestBytesToNibbles(t *testing.T) {
	testCases := []struct {
		name    string
		bytes   []byte
		nibbles []byte
	}{
		{
			name:    "Example #1",
			bytes:   []byte{0x04, 0x12, 0x34, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			nibbles: []byte{0x0, 0x4, 0x1, 0x2, 0x3, 0x4, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF},
		},

		{
			name:    "Example #2",
			bytes:   []byte{0x12, 0x34, 0x50},
			nibbles: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x0},
		},

		{
			name:    "Example #3",
			bytes:   []byte{0x78, 0x09, 0xBE},
			nibbles: []byte{0x7, 0x8, 0x0, 0x9, 0xB, 0xE},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := BytesToNibbles(tc.bytes)
			test.AssertBytesEqual(t, tc.nibbles, result)
		})
	}
}

func TestNibblesToBytes(t *testing.T) {
	testCases := []struct {
		name    string
		nibbles []byte
		bytes   []byte
	}{
		{
			name:    "Example #1",
			nibbles: []byte{0x0, 0x4, 0x1, 0x2, 0x3, 0x4, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF},
			bytes:   []byte{0x04, 0x12, 0x34, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},

		{
			name:    "Odd number of Nibbles",
			nibbles: []byte{0x1, 0x2, 0x3, 0x4, 0x5},
			bytes:   []byte{0x12, 0x34, 0x50},
		},

		{
			name:    "Even number of Nibbles",
			nibbles: []byte{0x7, 0x8, 0x0, 0x9, 0xB, 0xE},
			bytes:   []byte{0x78, 0x09, 0xBE},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NibblesToBytes(tc.nibbles)
			test.AssertBytesEqual(t, tc.bytes, result)
		})
	}
}

func TestParseHexNibbles(t *testing.T) {
	t.Run("Some examples", func(t *testing.T) {
		testCases := []struct {
			name    string
			text    string
			nibbles []byte
		}{
			{
				name:    "Odd length example",
				text:    "12345",
				nibbles: []byte{0x1, 0x2, 0x3, 0x4, 0x5},
			},
			{
				name:    "Even length example",
				text:    "7809BE",
				nibbles: []byte{0x7, 0x8, 0x0, 0x9, 0xB, 0xE},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := ParseHexNibbles(tc.text)
				require.NoError(t, err)

				test.AssertBytesEqual(t, tc.nibbles, result)
			})
		}
	})

	t.Run("Compare to hex.DecodeString", func(t *testing.T) {
		testCases := []struct {
			name string
			text string
		}{
			{
				name: "999000428685",
				text: "999000428685",
			},
			{
				name: "049589FFFFFFFFFF",
				text: "049589FFFFFFFFFF",
			},
			{
				name: "0000999000428685",
				text: "0000999000428685",
			},
			{
				name: "0495106FFFBD797A",
				text: "0495106FFFBD797A",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				bytesFromHex := test.MustParseHex(t, tc.text)

				nibbles, err := ParseHexNibbles(tc.text)
				require.NoError(t, err)

				bytesFromNibbles := NibblesToBytes(nibbles)

				test.AssertBytesEqual(t, bytesFromHex, bytesFromNibbles)
			})
		}
	})
}

func TestTruncateLeft(t *testing.T) {
	t.Run("Manual Examples", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  string
			length int
			output string
		}{
			{
				name:   "Specified length shorter than length #1",
				input:  "1234567890",
				length: 5,
				output: "67890",
			},
			{
				name:   "Specified length shorter than length #2",
				input:  "abcdefghij",
				length: 9,
				output: "bcdefghij",
			},
			{
				name:   "Specified length equal to length",
				input:  "1234567890",
				length: 10,
				output: "1234567890",
			},
			{
				name:   "Specified length longer than length",
				input:  "1234567890",
				length: 11,
				output: "1234567890",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := TruncateLeft([]byte(tc.input), tc.length)
				test.AssertBytesEqual(t, []byte(tc.output), result)
			})
		}
	})
}

func TestTruncateRight(t *testing.T) {
	t.Run("Manual Examples", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  string
			length int
			output string
		}{
			{
				name:   "Specified length shorter than length #1",
				input:  "1234567890",
				length: 5,
				output: "12345",
			},
			{
				name:   "Specified length shorter than length #2",
				input:  "abcdefghij",
				length: 9,
				output: "abcdefghi",
			},
			{
				name:   "Specified length equal to length",
				input:  "1234567890",
				length: 10,
				output: "1234567890",
			},
			{
				name:   "Specified length longer than length",
				input:  "1234567890",
				length: 11,
				output: "1234567890",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := TruncateRight([]byte(tc.input), tc.length)
				test.AssertBytesEqual(t, []byte(tc.output), result)
			})
		}
	})
}
