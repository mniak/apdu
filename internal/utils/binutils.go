package utils

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"strings"
)

// XOR performs an Exclusive Or operation on the data and returns a slice with the results.
// This function does not mutate the slices fed into it.
//
// If the length of `other` is smaller than the length of `a`, it will be padded to the left, adding zeroes to the left of b.
// If the length of `other` is greater than the length of `a`, it will be truncated to the left, removing the most significant bits.
func XOR(a []byte, others ...[]byte) []byte {
	lenA := len(a)
	result := make([]byte, lenA)
	copy(result, a)
	for _, n := range others {
		n = PadLeft(n, 0x00, lenA)[:lenA]

		for i := 0; i < lenA; i++ {
			result[i] ^= n[i]
		}
	}

	return result
}

// XORInPlace performs an Exclusive Or operation on the data mutating it.
//
// If the length of `other` is smaller than the length of `data`, it will be truncated to the left, removing the most significant bits.
// If the length of `other` is greater than the length of `data`, it will be padded to the left, adding zeroes to the left.
func XORInPlace(data, other []byte) {
	other = PadLeft(other, 0x00, len(data))[:len(data)]

	for i := 0; i < len(data); i++ {
		data[i] = data[i] ^ other[i]
	}
}

func InvertBits(data []byte) []byte {
	newData := make([]byte, len(data))
	for i, d := range data {
		newData[i] = d ^ 0xFF
	}
	return newData
}

func PadLeft(data []byte, padbyte byte, totalLength int) []byte {
	if totalLength <= len(data) {
		return data[len(data)-totalLength:]
	}
	padding := bytes.Repeat([]byte{padbyte}, totalLength-len(data))
	return append(padding, data...)
}

func PadRight(data []byte, padbyte byte, totalLength int) []byte {
	if totalLength <= len(data) {
		return data[:totalLength]
	}
	padding := bytes.Repeat([]byte{padbyte}, totalLength-len(data))
	return append(data, padding...)
}

// TruncateLeft discards bytes at the left until the length is equal the specified.
// If the length of the data is already equal to or shorter than `totalLength`, the
// data will be returned unchanged.
func TruncateLeft(data []byte, totalLength int) []byte {
	if len(data) <= totalLength {
		return data
	}
	return data[len(data)-totalLength:]
}

// TruncateRight discards bytes at the right until the length is equal the
// specified. If the length of the data is already equal to or shorter than
// `totalLength`, the data will be returned unchanged.
func TruncateRight(data []byte, totalLength int) []byte {
	if len(data) <= totalLength {
		return data
	}
	return data[:totalLength]
}

// Pad80 pads the data with byte 0x80 and then 0x00s until the data is an exact
// multiple of the block size.
//
// `skipIfExact` controls if an entire block should be added when the data is an
// exact multiple of the block size or if the padding should be skipped entirely.
func Pad80(data []byte, skipIfExact bool) []byte {
	const blockSize = 8
	if len(data)%blockSize == 0 && skipIfExact {
		return data
	}

	// 4. add padding to data
	data = append(data, 0x80)
	for len(data)%blockSize > 0 {
		data = append(data, 0x00)
	}
	return data
}

func SetNibble(slice []byte, index int, nibble uint8) {
	if index >= len(slice)*2 {
		panic("invalid index")
	}
	nibble &= 0xF
	byteIndex := index / 2

	shift := (1 - (index % 2)) * 4
	positiveMask := 0x0F << shift

	slice[byteIndex] &^= byte(positiveMask)
	slice[byteIndex] |= nibble << uint8(shift)
}

func BytesToNibbles(byts []byte) []byte {
	var buffer bytes.Buffer
	for _, b := range byts {
		buffer.WriteByte(b >> 4)
		buffer.WriteByte(b & 0x0F)
	}
	return buffer.Bytes()
}

func NibblesToBytes(nibbles []byte) []byte {
	slice := make([]byte, (len(nibbles)+1)/2)

	for idx, n := range nibbles {
		bidx := idx / 2

		shift := (1 - (idx % 2)) * 4
		slice[bidx] |= (n & 0xF) << uint8(shift)
	}
	return slice
}

func ParseHexNibbles(src string) ([]byte, error) {
	dst := make([]byte, len(src))
	for i, char := range src {
		var nibble byte
		if char >= '0' && char <= '9' {
			nibble = byte(char) - byte('0')
		} else if char >= 'A' && char <= 'F' {
			nibble = byte(char) - byte('A') + 10
		} else if char >= 'a' && char <= 'f' {
			nibble = byte(char) - byte('a') + 10
		}
		if nibble > 0x0f {
			return nil, hex.InvalidByteError(char)
		}
		dst[i] = nibble
	}
	return dst, nil
}

func NibblesToHex(nibbles []byte) string {
	var sb strings.Builder
	for _, n := range nibbles {
		sb.WriteString(strconv.Itoa(int(n) % 10))
	}
	return sb.String()
}
