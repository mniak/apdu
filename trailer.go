package apdu

import "fmt"

type Trailer uint16

func NewTrailer(sw1, sw2 byte) Trailer {
	return Trailer(sw1)<<8 | Trailer(sw2)
}

func (t Trailer) SW1() byte {
	return byte(t >> 8)
}

func (t Trailer) SW2() byte {
	return byte(t)
}

func (t Trailer) message() string {
	switch t {
	case 0x9000:
		return "command normally completed"

	case 0x6E00:
		return "CLA not supported"
	case 0x6D00:
		return "CLA supported, but INS not programmed or invalid"
	case 0x6B00:
		return "CLA INS supported, but P1 P2 incorrect"
	case 0x6700:
		return "CLA INS P1 P2 supported, but P3 incorrect"
	case 0x6F00:
		return "command not supported and no precise diagnosis given"

	case 0x6A00:
		return "no information given"
	case 0x6A80:
		return "incorrect parameters in the command data field"
	case 0x6A81:
		return "function not supported"
	case 0x6A82:
		return "file or application not found"
	case 0x6A83:
		return "record not found"
	case 0x6A84:
		return "not enough memory space in the file"
	case 0x6A85:
		return "Nc inconsistent with TLV structure"
	case 0x6A86:
		return "incorrect parameters P1-P2"
	case 0x6A87:
		return "Nc inconsistent with parameters P1-P2"
	case 0x6A88:
		return "referenced data or reference data not found (exact meaning depending on the command)"
	case 0x6A89:
		return "file already exists"
	case 0x6A8A:
		return "DF name already exists"

	}
	return ""
}

func (t Trailer) code() string {
	return fmt.Sprintf("%02X %02X", t.SW1(), t.SW2())
}

func (t Trailer) String() string {
	msg := t.message()
	if msg != "" {
		return msg
	}
	return t.code()
}
