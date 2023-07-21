package apdu

import (
	"fmt"
	"testing"

	"github.com/mniak/apdu/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestClass_Invalid(t *testing.T) {
	template := test.ByteTemplate("1111_1111")
	t.Run("FF is invalid", func(t *testing.T) {
		class := Class(template.Random(t))
		assert.True(t, class.Invalid())
	})

	t.Run("Any other value is valid", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			randomValue := template.BadValue(t)
			t.Run(fmt.Sprintf("%08b", randomValue), func(t *testing.T) {
				class := Class(randomValue)
				assert.False(t, class.Invalid())
			})
		}
	})
}

// func TestClass_FirstInterindustry(t *testing.T) {
// 	template := test.ByteTemplate("000x_xxxx")

// 	t.Run("Yes", func(t *testing.T) {
// 		for i := 0; i < 10; i++ {
// 			randomValue := template.Random(t)
// 			t.Run(fmt.Sprintf("%08b", randomValue), func(t *testing.T) {
// 				class := Class(randomValue)
// 				assert.True(t, class.FirstInterindustry())
// 			})
// 		}
// 	})
// 	t.Run("No", func(t *testing.T) {
// 		for i := 0; i < 10; i++ {
// 			randomValue := template.BadValue(t)
// 			t.Run(fmt.Sprintf("%08b", randomValue), func(t *testing.T) {
// 				class := Class(randomValue)
// 				assert.False(t, class.FirstInterindustry())
// 			})
// 		}
// 	})
// }

func TestClass_Proprietary(t *testing.T) {
	testdata := []struct {
		template test.ByteTemplate
		expected bool
	}{
		{
			template: "0xxx_xxxx",
			expected: false,
		},
		{
			template: "1xxx_xxxx",
			expected: true,
		},
	}
	for _, td := range testdata {
		t.Run(fmt.Sprintf("%s,proprietary=%v", td.template.String(), td.expected), func(t *testing.T) {
			t.Run(fmt.Sprintf("Min:%08b", td.template.Min(t)), func(t *testing.T) {
				class := Class(td.template.Min(t))
				assert.Equal(t, td.expected, class.Proprietary())
			})

			t.Run(fmt.Sprintf("Max:%08b", td.template.Max(t)), func(t *testing.T) {
				class := Class(td.template.Max(t))
				assert.Equal(t, td.expected, class.Proprietary())
			})

			for i := 0; i < 10; i++ {
				randomValue := td.template.Random(t)
				t.Run(fmt.Sprintf("%08b", randomValue), func(t *testing.T) {
					class := Class(randomValue)
					assert.Equal(t, td.expected, class.Proprietary())
				})
			}
		})
	}
}

func TestClass_LastInChain(t *testing.T) {
	testdata := []struct {
		template test.ByteTemplate
		expected bool
	}{
		// First interindustry values
		{
			template: "0000_xxxx",
			expected: true,
		},
		{
			template: "0001_xxxx",
			expected: false,
		},

		// Further interindustry values
		{
			template: "01x0_xxxx",
			expected: true,
		},
		{
			template: "01x1_xxxx",
			expected: false,
		},

		// Unspecified values
		{
			template: "0010_xxxx",
			expected: true,
		},
		{
			template: "0011_xxxx",
			expected: false,
		},
		{
			template: "1xxx_xxxx",
			expected: false,
		},
	}
	for _, td := range testdata {
		t.Run(fmt.Sprintf("%s,last=%v", td.template.String(), td.expected), func(t *testing.T) {
			t.Run(fmt.Sprintf("Min:%08b", td.template.Min(t)), func(t *testing.T) {
				class := Class(td.template.Min(t))
				assert.Equal(t, td.expected, class.LastInChain())
			})

			t.Run(fmt.Sprintf("Max:%08b", td.template.Max(t)), func(t *testing.T) {
				class := Class(td.template.Max(t))
				assert.Equal(t, td.expected, class.LastInChain())
			})

			for i := 0; i < 10; i++ {
				randomValue := td.template.Random(t)
				t.Run(fmt.Sprintf("%08b", randomValue), func(t *testing.T) {
					class := Class(randomValue)
					assert.Equal(t, td.expected, class.LastInChain())
				})
			}
		})
	}
}

// func TestClass_SecureMessaging(t *testing.T) {
// 	testdata := []struct {
// 		template test.ByteTemplate
// 		expected SecureMessagingIndication
// 	}{
// 		// First interindustry values
// 		{
// 			template: "000x_00xx",
// 			expected: NoSMOrNoIndication,
// 		},
// 		{
// 			template: "000x_01xx",
// 			expected: ProprietarySMFormat,
// 		},
// 		{
// 			template: "000x_10xx",
// 			expected: SMAccordingToSection6_NotProcessed,
// 		},
// 		{
// 			template: "000x_11xx",
// 			expected: SMAccordingToSection6_Authenticated,
// 		},

// 		// Futher interindustry values
// 		{
// 			template: "010x_00xx",
// 			expected: NoSMOrNoIndication,
// 		},
// 		{
// 			template: "011x_00xx",
// 			expected: SMAccordingToSection6_NotProcessed,
// 		},

// 		// Unspecified values
// 		{
// 			template: "001x_11xx",
// 			expected: NoSMOrNoIndication,
// 		},
// 		{
// 			template: "1xxx_11xx",
// 			expected: NoSMOrNoIndication,
// 		},
// 	}
// 	for _, td := range testdata {
// 		t.Run(fmt.Sprintf("%s,last=%v", td.template.String(), td.expected), func(t *testing.T) {
// 			t.Run(fmt.Sprintf("Min:%08b", td.template.Min(t)), func(t *testing.T) {
// 				class := Class(td.template.Min(t))
// 				assert.Equal(t, td.expected, class.SecureMessagingIndication())
// 			})

// 			t.Run(fmt.Sprintf("Max:%08b", td.template.Max(t)), func(t *testing.T) {
// 				class := Class(td.template.Max(t))
// 				assert.Equal(t, td.expected, class.SecureMessagingIndication())
// 			})

// 			for i := 0; i < 10; i++ {
// 				randomValue := td.template.Random(t)
// 				t.Run(fmt.Sprintf("%08b", randomValue), func(t *testing.T) {
// 					class := Class(randomValue)
// 					assert.Equal(t, td.expected, class.SecureMessagingIndication())
// 				})
// 			}
// 		})
// 	}
// }

// func TestClass_Invalids(t *testing.T) {
// 	testdata := []struct {
// 		name string
// 		min  byte
// 		max  byte
// 	}{
// 		{
// 			min:  "0010_0000",
// 			max:  "0011_1111",
// 		},
// 	}
// }
