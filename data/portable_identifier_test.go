package data

import (
	"testing"
	"time"
)

func TestMakeIdentifier_Zero(t *testing.T) {
	id := MakeIdentifier(time.Unix(0, 0), 0)
	if id.String() != "aaaaaaa-aaaaaaa" {
		t.Errorf("expected 'aaaaaaa-aaaaaaa' but got %s", id.String())
	}

	id = MakeIdentifier(time.Unix(0b0001_1000_1000_0010, 0), 0b0010_0000_1100_0100_0001)
	if id.String() != "dcbaaaa-aaaedcb" {
		t.Errorf("expected 'dcbaaaa-aaaedcb' but got %s", id.String())
	}

	oldBytes := randomBytes
	randomBytes = [14]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	id = MakeIdentifier(time.Unix(0b1110_0111_0111_1101, 0), 0b1111_1111_1111_1111_1111_1111_1111_1111_1111_1101_1111_0011_1011_1110)
	if id.String() != "dcbaaaa-aaaedcb" {
		t.Errorf("expected 'dcbaaaa-aaaedcb' but got %s", id.String())
	}
	randomBytes = oldBytes
}

func TestIdentifier_String(t *testing.T) {
	id := Identifier([]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13})
	if id.String() != "abcdefg-hjkmnpr" {
		t.Errorf("Expected abcdefg-hjkmnpr but got %s", id.String())
	}

	id = Identifier([]uint8{14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27})
	if id.String() != "stABCDE-FGHJKMN" {
		t.Errorf("Expected stABCDE-FGHJKMN but got %s", id.String())
	}
}

func TestPackedIdentifier_PackAndUnpack(t *testing.T) {
	id := MakeIdentifier(time.Unix(0b0001_1000_1000_0010, 0), 0b0010_0000_1100_0100_0001)

	packedId := id.Pack()
	id = packedId.Unpack()

	if id.String() != "dcbaaaa-aaaedcb" {
		t.Errorf("expected 'dcbaaaa-aaaedcb' but got %s", id.String())
	}
}
