package data

import (
	"errors"
	"fmt"
	"time"
)

type Identifier []uint8
type PackedIdentifier []uint8

var encodingSymbols = []uint8{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'm', 'n', 'p', 'r', 's', 't',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'M', 'N', 'P', 'R', 'S', 'T',
}

var encodingKey = map[byte]uint8{
	'a': 0, 'b': 1, 'c': 2, 'd': 3, 'e': 4, 'f': 5, 'g': 6, 'h': 7,
	'j': 8, 'k': 9, 'm': 10, 'n': 11, 'p': 12, 'r': 13, 's': 14, 't': 15,
	'A': 16, 'B': 17, 'C': 18, 'D': 19, 'E': 20, 'F': 21, 'G': 22, 'H': 23,
	'J': 24, 'K': 25, 'M': 26, 'N': 27, 'P': 28, 'R': 29, 'S': 30, 'T': 31,
}

var mask uint8 = 0x1F
var randomBytes = [14]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func MakeIdentifier(time time.Time, id int64) Identifier {
	timePortion := uint64(time.Unix())
	uId := uint64(id)
	return []uint8{
		(uint8(timePortion>>11) & mask) ^ randomBytes[0],
		(uint8(timePortion>>6) & mask) ^ randomBytes[1],
		(uint8(timePortion>>1) & mask) ^ randomBytes[2],
		(uint8(timePortion<<4)&0x10 | uint8(uId>>50)&0x0F) ^ randomBytes[3],
		(uint8(uId>>45) & mask) ^ randomBytes[4],
		(uint8(uId>>40) & mask) ^ randomBytes[5],
		(uint8(uId>>35) & mask) ^ randomBytes[6],
		(uint8(uId>>30) & mask) ^ randomBytes[7],
		(uint8(uId>>25) & mask) ^ randomBytes[8],
		(uint8(uId>>20) & mask) ^ randomBytes[9],
		(uint8(uId>>15) & mask) ^ randomBytes[10],
		(uint8(uId>>10) & mask) ^ randomBytes[11],
		(uint8(uId>>5) & mask) ^ randomBytes[12],
		(uint8(uId) & mask) ^ randomBytes[13],
	}
}

func ParseIdentifier(id string) (Identifier, error) {
	if len(id) != 15 && len(id) != 14 {
		return Identifier{}, errors.New("parse exception - identifier is 15 characters or 14 characters without the '-'")
	}

	result := make([]uint8, 14)
	idx := 0
	for _, r := range id {
		if r == '-' {
			continue
		}
		v, ok := encodingKey[uint8(r)]
		if !ok {
			return Identifier{}, errors.New("invalid character in identifier string")
		}

		result[idx] = v
		idx++
	}

	return result, nil
}

func (i Identifier) String() string {
	front := []byte{
		encodingSymbols[i[0]&mask],
		encodingSymbols[i[1]&mask],
		encodingSymbols[i[2]&mask],
		encodingSymbols[i[3]&mask],
		encodingSymbols[i[4]&mask],
		encodingSymbols[i[5]&mask],
		encodingSymbols[i[6]&mask],
	}
	back := []byte{
		encodingSymbols[i[7]&mask],
		encodingSymbols[i[8]&mask],
		encodingSymbols[i[9]&mask],
		encodingSymbols[i[10]&mask],
		encodingSymbols[i[11]&mask],
		encodingSymbols[i[12]&mask],
		encodingSymbols[i[13]&mask],
	}

	return fmt.Sprintf("%s-%s", string(front), string(back))
}

func (i Identifier) Pack() PackedIdentifier {
	result := make([]uint8, 9)
	srcMask := []uint8{
		0b0_0001, // 0  1
		0b0_0011, // 1  2
		0b0_0111, // 2  3
		0b0_1111, // 3  4
		0b1_1111, // 4  5 0 offset bitsPacked % 5 is 0  5 or more bits available
		0b1_1110, // 5  4 1 offset bitsPacked % 5 is 1  4 bits available
		0b1_1100, // 6  3 2 offset bitsPacked % 5 is 2  3 bits available
		0b1_1000, // 7  2 offset bitsPacked % 5 is 3  2 bits available
		0b1_0000, // 8  1 offset bitsPacked % 5 is 4  1 bit available
	}

	bitsPacked := 0
	for bitsPacked < 70 {
		sourceByte := bitsPacked / 5
		destByte := bitsPacked / 8
		nextByte := destByte + 1
		bitsAvailable := 8 - (bitsPacked % 8)

		if bitsAvailable > 4 {
			// If there are 5 or more bits available in the byte we're packing into, then
			// shift bits left and add them to the destination byte.
			result[destByte] = result[destByte] | (i[sourceByte] << (bitsAvailable - 5))
		} else {
			// If there are less than 5 bits available, then the 5-bit value is smeared
			// across two 8-bit bytes.  The high bits taken from the front of the 5-bit value
			// are the masks at indexes 8, 7, 6, 5, for 1, 2, 3, or 4 bits available.  The
			// low masks are 3, 2, 1, and 0 for 1, 2, 3, or 4 bits available.  The values have
			// to be shifted down for the remainder of the current byte and shifted up for the
			// start of the next byte.
			highMask := srcMask[9-bitsAvailable]
			lowMask := srcMask[4-bitsAvailable]
			result[destByte] = result[destByte] | ((i[sourceByte] & highMask) >> (5 - bitsAvailable))
			result[nextByte] = result[nextByte] | ((i[sourceByte] & lowMask) << (3 + bitsAvailable))
		}
		bitsPacked += 5
	}

	return result
}

func (p PackedIdentifier) Unpack() Identifier {
	result := make([]uint8, 14)

	masks := []uint8{
		0b0000_0001, //  0                       bitsPacked % 8 was 4
		0b0000_0011, //  1                       bitsPacked % 8 was 5
		0b0000_0111, //  2                       bitsPacked % 8 was 6
		0b0000_1111, //  3                       bitsPacked % 8 was 7
		0b0001_1111, //  4 bitsPacked % 8 is 0
		0b0011_1110, //  5 bitsPacked % 8 is 1
		0b0111_1100, //  6 bitsPacked % 8 is 2
		0b1111_1000, //  7 bitsPacked % 8 is 3
		0b1111_0000, //  8 bitsPacked % 8 is 4
		0b1110_0000, //  9 bitsPacked % 8 is 5
		0b1100_0000, // 10 bitsPacked % 8 is 6
		0b1000_0000, // 11 bitsPacked % 8 is 7
	}

	bitsUnpacked := 0
	for bitsUnpacked < 70 {
		sourceByte := bitsUnpacked / 8
		nextByte := sourceByte + 1
		destByte := uint8(0)

		remainder := bitsUnpacked % 8
		if remainder < 4 {
			// In this case there are 0, 1, 2, or 3 bits already unpacked in the current
			// byte.  So the mask is 7, 6, 5, 4 for 0, 1, 2, 3 bits unpacked, respectively.
			destByte = (p[sourceByte] & masks[7-remainder]) >> (3 - remainder)
		} else {
			// In this case the 5 bit value is smeared across 2 8-bit bytes.  The current
			// byte holds the high bits and the next byte holds the low bits.  The mask
			// for the high bits in the current byte will be 3, 2, 1, or 0, if there are
			// 4, 3, 2, or 1 bits remaining.  The low mask will be the 11, 10, 9, 8 if
			// there are 4, 3, 2, or 1 bits remaining.  The values are shifted left to move
			// the high bits into place and shifted right to move the low bits into place.
			// The high bits come from the current 8-bit byte processed and the low come
			// from the next 8-bit byte.
			highMask := masks[7-remainder]
			lowMask := masks[12-(remainder-3)]
			destByte = (p[sourceByte] & highMask) << (remainder - 3)
			destByte += (p[nextByte] & lowMask) >> (11 - remainder)
		}
		result[bitsUnpacked/5] = destByte
		bitsUnpacked += 5
	}

	return result
}
