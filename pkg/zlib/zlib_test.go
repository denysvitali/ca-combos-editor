package zlib_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/denysvitali/ca-combos-editor/pkg/zlib"
)

func TestRoundTrip(t *testing.T) {
	original, err := os.ReadFile("../../test/resources/2019-10-17/extracted")
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}

	compressed, err := zlib.Compress(bytes.NewReader(original))
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	decompressed, err := zlib.Decompress(bytes.NewReader(compressed))
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(decompressed, original) {
		t.Fatalf("round-trip mismatch: got %d bytes, want %d bytes", len(decompressed), len(original))
	}
}

func TestDecompressOriginal(t *testing.T) {
	want, err := os.ReadFile("../../test/resources/2019-10-17/extracted")
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}

	f, err := os.Open("../../test/resources/2019-10-17/00028874")
	if err != nil {
		t.Fatalf("failed to open 00028874: %v", err)
	}
	defer func() { _ = f.Close() }()

	got, err := zlib.Decompress(f)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("decompressed mismatch: got %d bytes, want %d bytes", len(got), len(want))
	}
}
