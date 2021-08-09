package tcg

import "testing"

func TestEncode(t *testing.T) {
	encoded := fpeEncode(`(\d\d)-(\d\d)`, "12-34")
	if encoded != "56-78" {
		t.Fail()
	}
}

func TestDecode(t *testing.T) {
	encoded := fpeDecode(`(\d\d)-(\d\d)`, "56-78")
	if encoded != "12-34" {
		t.Fail()
	}
}
