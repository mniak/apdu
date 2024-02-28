package apdu

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseResponse_Examples(t *testing.T) {
	testdata := []struct {
		name     string
		hexdata  string
		expected Response
	}{
		{
			name:    "SELECT response",
			hexdata: "6153",
			expected: Response{
				Data:    []byte{},
				Trailer: 0x6153,
			},
		},

		// {
		// 	name: "Select",
		// 	data: "6F518407A0000000031010A546500C56495341204352454449544F8701029F120C56495341204352454449544F5F2D0870746573656E66729F1101019F38039F1A02BF0C0E9F5A0652098600763042034901449000",
		// },
	}
	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			data, err := hex.DecodeString(td.hexdata)
			require.NoError(t, err)

			resp, err := ParseResponse(data)
			assert.NoError(t, err)
			assert.Equal(t, td.expected, resp)
		})
	}
}
