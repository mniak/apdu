package test

import (
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

type ByteTemplate string

func (bt ByteTemplate) Random(t *testing.T) byte {
	t.Helper()
	template := bt.checkTemplateLength(t)

	var result byte
	for _, ch := range template {
		result <<= 1

		switch ch {
		case '0':
		case '1':
			result += 1
		case 'x', 'X':
			result += byte(gofakeit.IntRange(0, 1))
		default:
			bt.failTemplateChar(t, ch)
		}
	}
	return result
}

func (bt ByteTemplate) BadValue(t *testing.T) byte {
	t.Helper()
	template := bt.checkTemplateLength(t)

	explicitBits := make([]int, 0)
	var result byte
	for i, ch := range template {
		result <<= 1

		switch ch {
		case '0':
			explicitBits = append(explicitBits, int(i))
		case '1':
			explicitBits = append(explicitBits, int(i))
			result += 1
		case 'x', 'X':
			result += byte(gofakeit.IntRange(0, 1))
		default:
			bt.failTemplateChar(t, ch)
		}
	}

	if len(explicitBits) == 0 {
		t.Log("None of the bits are explicitly set on the template. At least on character must not be 'x'")
	}
	gofakeit.ShuffleInts(explicitBits)
	badBitAmount := gofakeit.IntRange(1, len(explicitBits))
	for _, badBit := range explicitBits[:badBitAmount] {
		badBitMask := 1 << badBit
		result ^= byte(badBitMask)
	}
	return result
}

func (bt ByteTemplate) Min(t *testing.T) byte {
	t.Helper()
	template := bt.checkTemplateLength(t)

	var result byte
	for _, ch := range template {
		result <<= 1

		switch ch {
		case '0', 'x', 'X':
		case '1':
			result += 1
		default:
			bt.failTemplateChar(t, ch)
		}
	}
	return result
}

func (bt ByteTemplate) Max(t *testing.T) byte {
	t.Helper()
	template := bt.checkTemplateLength(t)

	var result byte
	for _, ch := range template {
		result <<= 1

		switch ch {
		case '0':
		case '1', 'x', 'X':
			result += 1
		default:
			bt.failTemplateChar(t, ch)
		}
	}
	return result
}

func (bt ByteTemplate) failTemplateChar(t *testing.T, ch rune) {
	t.Logf("Invalid byte template '%s'. Characters must be 1, 0 or x. Found char '%v'", string(bt), ch)
}

func (bt ByteTemplate) String() string {
	return strings.ReplaceAll(string(bt), "_", "")
}

func (bt ByteTemplate) checkTemplateLength(t *testing.T) string {
	template := bt.String()
	require.Lenf(t, template, 8, "Invalid byte template '%s'", template)
	return template
}

func (bt ByteTemplate) Validate(t *testing.T, b byte) bool {
	template := bt.checkTemplateLength(t)
	for i := 0; i < 8; i++ {
		ch := template[i]
		bit := b >> (8 - i - 1) & 1

		switch ch {
		case 'x', 'X':
			continue
		case '0':
			if bit != 0 {
				return false
			}
		case '1':
			if bit != 1 {
				return false
			}
		default:
			bt.failTemplateChar(t, rune(ch))
		}
	}
	return true
}
