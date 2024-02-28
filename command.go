package apdu

import (
	"bytes"
	"fmt"

	"github.com/mniak/tlv"
	"github.com/pkg/errors"
)

type Command struct {
	Class            Class
	Instruction      Instruction
	Parameters       Parameters
	Data             []byte
	MaxReponseLength byte
}
type Parameters struct {
	P1 byte
	P2 byte
}

func (c Command) Bytes(enc tlv.LengthEncoder) ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte(byte(c.Class))
	b.WriteByte(byte(c.Instruction))
	b.WriteByte(byte(c.Parameters.P1))
	b.WriteByte(byte(c.Parameters.P2))

	if len(c.Data) > 0 {
		lengthBytes, err := enc.Encode(len(c.Data))
		if err != nil {
			return nil, errors.WithMessage(err, "failed to encode command length")
		}
		b.Write(lengthBytes)
		b.Write(c.Data)
	}

	maxResponseLengthBytes, err := enc.Encode(int(c.MaxReponseLength))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to encode maximum expected response length")
	}
	b.Write(maxResponseLengthBytes)
	return b.Bytes(), nil
}

func (c Command) String() string {
	return fmt.Sprintf("CLA=%02X INS=%02X P1=%02X P2=%02X DATA=[%2X] Le=%02X", c.Class, c.Instruction, c.Parameters.P1, c.Parameters.P2, c.Data, c.MaxReponseLength)
}

func (c Command) StringPretty() string {
	return fmt.Sprintf("Class: %02X\nInstruction: %02X\nP1: %02X\nP2: %02X\nData: [%2X]\nLe: %02X",
		c.Class,
		c.Instruction,
		c.Parameters.P1, c.Parameters.P2,
		c.Data,
		c.MaxReponseLength,
	)
}
