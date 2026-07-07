package readers

import (
	"bytes"
	"fmt"
	"io"
)

// BinaryReader wraps a bytes.Reader and provides helpers for reading the
// little-endian Qualcomm NV item format.
type BinaryReader struct {
	reader *bytes.Reader
}

// NewBinaryReader creates a BinaryReader backed by the provided bytes.Reader.
func NewBinaryReader(reader *bytes.Reader) BinaryReader {
	return BinaryReader{reader: reader}
}

// NewMyReader is a deprecated alias for NewBinaryReader.
//
// Deprecated: use NewBinaryReader instead.
func NewMyReader(reader *bytes.Reader) BinaryReader {
	return NewBinaryReader(reader)
}

// ReadByte reads the next byte or returns an error on EOF.
func (m *BinaryReader) ReadByte() (byte, error) {
	b, err := m.reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read byte: %w", err)
	}
	return b, nil
}

// Rb reads the next byte or returns an error on EOF.
//
// Deprecated: use ReadByte instead.
func (m *BinaryReader) Rb() (byte, error) {
	return m.ReadByte()
}

// ReadBytes reads exactly num bytes.
func (m *BinaryReader) ReadBytes(num int) ([]byte, error) {
	bArr := make([]byte, num)
	if _, err := io.ReadFull(m.reader, bArr); err != nil {
		return nil, fmt.Errorf("read %d bytes: %w", num, err)
	}
	return bArr, nil
}

// Expect reads a byte and verifies it matches the expected value.
func (m *BinaryReader) Expect(b byte) error {
	found, err := m.ReadByte()
	if err != nil {
		return err
	}
	if found != b {
		return fmt.Errorf("unexpected byte 0x%02X, expected 0x%02X", found, b)
	}
	return nil
}

// Len returns the number of bytes remaining in the reader.
func (m *BinaryReader) Len() int {
	return m.reader.Len()
}
