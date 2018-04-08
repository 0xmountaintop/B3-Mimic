package testutil

import (
	"bytes"
	"io"
	"testing"

	"github.com/bytom/protocol/bc"
)

func MustDecodeHash(s string) (h bc.Hash) {
	if err := h.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return h
}

func MustDecodeAsset(s string) (h bc.AssetID) {
	if err := h.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return h
}

func Serialize(t *testing.T, wt io.WriterTo) []byte {
	var b bytes.Buffer
	if _, err := wt.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
