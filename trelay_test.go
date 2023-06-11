package trelay_test

import (
	"bytes"
	"testing"

	"github.com/btvoidx/trelay"
)

func compbytes(t *testing.T, b1, b2 []byte) {
	if len(b1) != len(b2) {
		t.Fatalf("length mismatch (len(b1) == %d; len(b2) == %d)", len(b1), len(b2))
	}

	for i := range b1 {
		if b1[i] != b2[i] {
			t.Fatalf("value mismatch (b1[%d] == %x; b2[%d] == %x)", i, b1[i], i, b2[i])
		}
	}
}

func TestFscanSimple(t *testing.T) {
	// 15 0 1 "Terraria123"
	data := []byte{0xf, 0x0, 0x1, 0xb, 0x54, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x61, 0x31, 0x32, 0x33}

	var ln uint16
	var id byte
	var ver string

	_, err := trelay.Fscan(bytes.NewReader(data), &ln, &id, &ver)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	switch {
	case ln != 15:
		t.Fatalf("bad read: ln == %d; expected %d", ln, 15)
	case id != 1:
		t.Fatalf("bad read: id == %d; expected %d", id, 1)
	case ver != "Terraria123":
		t.Fatalf("bad read: ver == %q; expected %q", ver, "Terraria123")
	}
}

func TestFprintBuilder(t *testing.T) {
	p := &trelay.Packet{ID: 1}
	trelay.Fprint(p, "Terraria123")

	// 15 0 1 "Terraria123"
	data := []byte{0xf, 0x0, 0x1, 0xb, 0x54, 0x65, 0x72, 0x72, 0x61, 0x72, 0x69, 0x61, 0x31, 0x32, 0x33}

	compbytes(t, data, p.Bytes())
}
