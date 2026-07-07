package zlib

import (
	"bytes"
	"compress/zlib"
	"io"
)

// Compress reads all data from r and returns it compressed with zlib level 6.
func Compress(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, 6)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(w, r); err != nil {
		_ = w.Close()
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decompress reads a raw zlib stream from r and returns the decompressed bytes.
func Decompress(r io.Reader) ([]byte, error) {
	rc, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rc.Close() }()

	return io.ReadAll(rc)
}
