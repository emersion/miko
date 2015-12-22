package handler_test

import (
	"bytes"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
	"testing"
)

func TestRead_uint8(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x7})

	var val uint8
	err := handler.Read(b, &val)
	if err != nil {
		t.Fatal("Error while reading:", err)
	}
	if val != 7 {
		t.Fatal("Bad value:", val)
	}
}

func TestRead_uint16(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x3, 0xB})

	var val uint16
	err := handler.Read(b, &val)
	if err != nil {
		t.Fatal("Error while reading:", err)
	}
	if val != 779 {
		t.Fatal("Bad value:", val)
	}
}

func TestRead_int16(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x6, 0x4})

	var val int16
	err := handler.Read(b, &val)
	if err != nil {
		t.Fatal("Error while reading:", err)
	}
	if val != 1540 {
		t.Fatal("Bad value:", val)
	}
}

func TestRead_string(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x0, 0x3, 'a', 'b', 'c'})

	var val string
	err := handler.Read(b, &val)
	if err != nil {
		t.Fatal("Error while reading:", err)
	}
	if val != "abc" {
		t.Fatal("Bad value:", val)
	}
}

func TestRead_multiple(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x0, 0x1, 'a', 0x2, 0x0, 0xF, 0x0, 0x3, 'L', 'O', 'L'})

	var val1 string
	var val2 uint8
	var val3 uint16
	var val4 string
	err := handler.Read(b, &val1, &val2, &val3, &val4)
	if err != nil {
		t.Fatal("Error while reading:", err)
	}
	if val1 != "a" {
		t.Fatal("Bad value:", val1)
	}
	if val2 != 2 {
		t.Fatal("Bad value:", val2)
	}
	if val3 != 15 {
		t.Fatal("Bad value:", val3)
	}
	if val4 != "LOL" {
		t.Fatal("Bad value:", val4)
	}
}
