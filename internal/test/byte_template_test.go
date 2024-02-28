package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteTemplate_Validate(t *testing.T) {
	testdata := []struct {
		template ByteTemplate
		value    byte
		expected bool
	}{
		// 0s
		{
			template: "0000_0000",
			value:    0b0000_0000,
			expected: true,
		},
		{
			template: "0000_0000",
			value:    0b0000_0001,
			expected: false,
		},
		{
			template: "0000_0000",
			value:    0b0001_0000,
			expected: false,
		},
		// 1s
		{
			template: "1111_1111",
			value:    0b1111_1111,
			expected: true,
		},
		{
			template: "1111_1111",
			value:    0b1111_1110,
			expected: false,
		},
		{
			template: "1111_1111",
			value:    0b1110_1111,
			expected: false,
		},
		// x
		{
			template: "0000_x111",
			value:    0b0000_0111,
			expected: true,
		},
		{
			template: "0000_x111",
			value:    0b0000_0111,
			expected: true,
		},
		{
			template: "0000_x111",
			value:    0b0001_0111,
			expected: false,
		},
		{
			template: "0000_x111",
			value:    0b0001_0111,
			expected: false,
		},
	}
	for _, td := range testdata {
		t.Run(fmt.Sprintf("%s->%08b", td.template.String(), td.value), func(t *testing.T) {
			result := td.template.Validate(t, td.value)
			assert.Equal(t, td.expected, result)
		})
	}
}
