package readers

import (
	"bytes"
	"github.com/sirupsen/logrus"
)

type BinaryReader struct {
	reader *bytes.Reader
}

func NewMyReader(reader *bytes.Reader) BinaryReader {
	return BinaryReader{reader: reader}
}

func (m *BinaryReader) Rb() byte {
	b, e := m.reader.ReadByte()
	if e != nil {
		logrus.Fatal(e)
	}
	return b
}

func (m *BinaryReader) ReadBytes(num int) []byte {
	var bArr []byte
	for i := 0; i < num; i++ {
		bArr = append(bArr, m.Rb())
	}
	return bArr
}

func (m *BinaryReader) Expect(b byte) {
	found := m.Rb()
	if found != b {
		logrus.Fatalf("Unexpected byte %02X found, %02X expected",
			found&0xFF,
			b&0xFF)
	}
}
