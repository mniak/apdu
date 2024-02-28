package apdu

import (
	"encoding/hex"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/tlv"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestVerifyPlaintextPIN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeResponseBytes := []byte(gofakeit.SentenceSimple())

	mockClient := NewMockRawClient(ctrl)
	mockClient.EXPECT().
		SendCommand(gomock.Any()).
		Do(func(cmd Command) {
			assert.Equal(t, byte(0x00), byte(cmd.Class))
			assert.Equal(t, byte(0x20), byte(cmd.Instruction))
			assert.Equal(t, byte(0x00), byte(cmd.Parameters.P1))
			assert.Equal(t, byte(0x80), byte(cmd.Parameters.P2))
			assert.Equal(t, []byte{0x24, 0x12, 0x34, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, cmd.Data)

			cmdBytes, err := cmd.Bytes(tlv.ShortLengthEncoder)
			require.NoError(t, err)

			expectedBytes := lo.Must(hex.DecodeString("0020008008241234FFFFFFFFFF00"))
			assert.Equal(t, expectedBytes, cmdBytes)
		}).
		Return(Response{
			Data:    fakeResponseBytes,
			Trailer: NewTrailer(0x90, 0x00),
		}, nil)

	lowlevel := _LowLevelClient{
		RawClient: mockClient,
	}
	resp, err := lowlevel.VerifyPlaintextPIN([]int{1, 2, 3, 4})
	require.NoError(t, err)
	_ = resp
}
