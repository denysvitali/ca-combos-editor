package readers

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinaryReader_Rb(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    byte
		wantErr bool
	}{
		{"read first byte", []byte{0x01, 0x02}, 0x01, false},
		{"single byte input", []byte{0xAB}, 0xAB, false},
		{"empty reader", []byte{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMyReader(bytes.NewReader(tt.input))
			got, err := r.Rb()
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, io.EOF)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBinaryReader_Rb_EOF(t *testing.T) {
	r := NewMyReader(bytes.NewReader([]byte{0x01}))
	_, err := r.Rb()
	require.NoError(t, err)

	_, err = r.Rb()
	require.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestBinaryReader_ReadBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		n       int
		want    []byte
		wantErr bool
	}{
		{"read exact bytes", []byte{0x01, 0x02, 0x03}, 2, []byte{0x01, 0x02}, false},
		{"read all bytes", []byte{0xAA, 0xBB}, 2, []byte{0xAA, 0xBB}, false},
		{"read zero bytes", []byte{0x01}, 0, []byte{}, false},
		{"eof before enough bytes", []byte{0x01}, 2, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMyReader(bytes.NewReader(tt.input))
			got, err := r.ReadBytes(tt.n)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBinaryReader_Expect(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		expect  byte
		wantErr bool
	}{
		{"expected byte present", []byte{0x01, 0x02}, 0x01, false},
		{"unexpected byte", []byte{0x01}, 0x02, true},
		{"eof", []byte{}, 0x01, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMyReader(bytes.NewReader(tt.input))
			err := r.Expect(tt.expect)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestComboReader_ReadNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"simple number", "123", 123, false},
		{"number followed by letter", "42A", 42, false},
		{"number at end of string", "7", 7, false},
		{"empty input", "", -1, true},
		{"letter only", "A", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewComboReader(tt.input)
			got, err := r.ReadNumber()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComboReader_ReadClass(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"class A", "A", 1},
		{"class B", "B", 2},
		{"class C", "C", 3},
		{"class D", "D", 4},
		{"class E", "E", 5},
		{"invalid class", "Z", -1},
		{"empty input", "", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewComboReader(tt.input)
			got := r.ReadClass()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComboReader_HasNextBand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"has separator", "-1A", true},
		{"no separator", "1A", false},
		{"empty input", "", false},
		{"separator not first", "1A-", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewComboReader(tt.input)
			got := HasNextBand(&r)
			assert.Equal(t, tt.want, got)
		})
	}
}
