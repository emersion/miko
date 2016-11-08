package builder_test

import (
	"bytes"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"testing"
)

func bufferEquals(buf *bytes.Buffer, expected []byte) bool {
	got := buf.Bytes()

	if len(got) != len(expected) {
		return false
	}

	for i, b := range got {
		if b != expected[i] {
			return false
		}
	}

	return true
}

func TestWrite_uint8(t *testing.T) {
	b := new(bytes.Buffer)

	err := builder.Write(b, uint8(7))
	if err != nil {
		t.Fatal("Error while writing:", err)
	}
	if !bufferEquals(b, []byte{0x7}) {
		t.Fatal("Bad value:", b.Bytes())
	}
}

func TestWrite_uint16(t *testing.T) {
	b := new(bytes.Buffer)

	err := builder.Write(b, uint16(779))
	if err != nil {
		t.Fatal("Error while writing:", err)
	}
	if !bufferEquals(b, []byte{0x3, 0xB}) {
		t.Fatal("Bad value:", b.Bytes())
	}
}

func TestWrite_int16(t *testing.T) {
	b := new(bytes.Buffer)

	err := builder.Write(b, int16(1540))
	if err != nil {
		t.Fatal("Error while writing:", err)
	}
	if !bufferEquals(b, []byte{0x6, 0x4}) {
		t.Fatal("Bad value:", b.Bytes())
	}
}

func TestWrite_string(t *testing.T) {
	b := new(bytes.Buffer)

	err := builder.Write(b, "abc")
	if err != nil {
		t.Fatal("Error while writing:", err)
	}
	if !bufferEquals(b, []byte{0x0, 0x3, 'a', 'b', 'c'}) {
		t.Fatal("Bad value:", b.Bytes())
	}
}

func TestWrite_multiple(t *testing.T) {
	b := new(bytes.Buffer)

	err := builder.Write(b, "a", uint8(2), uint16(15), "LOL")
	if err != nil {
		t.Fatal("Error while writing:", err)
	}
	if !bufferEquals(b, []byte{0x0, 0x1, 'a', 0x2, 0x0, 0xF, 0x0, 0x3, 'L', 'O', 'L'}) {
		t.Fatal("Bad value:", b.Bytes())
	}
}
