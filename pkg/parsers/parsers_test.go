package parsers

import (
	"bytes"
	"testing"

	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func build13x(bands []types.Band) []byte {
	out := make([]byte, types.MaxBandsPerEntry*types.BandRecordSize13x)
	for i := 0; i < types.MaxBandsPerEntry; i++ {
		if i < len(bands) {
			out[i*3] = byte(bands[i].Band)
			out[i*3+1] = byte(bands[i].Band >> 8)
			out[i*3+2] = byte(bands[i].Class)
		}
	}
	return out
}

func build20x(bands []types.Band) []byte {
	out := make([]byte, types.MaxBandsPerEntry*types.BandRecordSize20x)
	for i := 0; i < types.MaxBandsPerEntry; i++ {
		if i < len(bands) {
			out[i*4] = byte(bands[i].Band)
			out[i*4+1] = byte(bands[i].Band >> 8)
			out[i*4+2] = byte(bands[i].Class)
			out[i*4+3] = byte(bands[i].Mimo)
		}
	}
	return out
}

func build33x(bands []types.Band) []byte {
	out := make([]byte, types.MaxBandsPerEntry*types.BandRecordSize33x)
	for i := 0; i < types.MaxBandsPerEntry; i++ {
		if i < len(bands) {
			out[i*11] = byte(bands[i].Band)
			out[i*11+1] = byte(bands[i].Band >> 8)
			out[i*11+2] = byte(bands[i].Class)
			for j := 0; j < types.AntennaCount; j++ {
				if j < len(bands[i].Antennas) {
					out[i*11+3+j] = byte(bands[i].Antennas[j])
				}
			}
		}
	}
	return out
}

func TestParse13xBands(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []types.Band
		wantErr  bool
	}{
		{
			name:     "single band",
			input:    build13x([]types.Band{{Band: 7, Class: 1}}),
			expected: []types.Band{{Band: 7, Class: 1}},
		},
		{
			name: "multiple bands",
			input: build13x([]types.Band{
				{Band: 7, Class: 1},
				{Band: 3, Class: 2},
				{Band: 1, Class: 3},
			}),
			expected: []types.Band{
				{Band: 7, Class: 1},
				{Band: 3, Class: 2},
				{Band: 1, Class: 3},
			},
		},
		{
			name: "empty slot skipped",
			input: build13x([]types.Band{
				{Band: 0, Class: 0},
				{Band: 7, Class: 1},
			}),
			expected: []types.Band{{Band: 7, Class: 1}},
		},
		{
			name: "invalid class skipped",
			input: build13x([]types.Band{
				{Band: 7, Class: 10},
				{Band: 3, Class: 1},
			}),
			expected: []types.Band{{Band: 3, Class: 1}},
		},
		{
			name:    "truncated input",
			input:   []byte{0x07, 0x00, 0x01},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := readers.NewBinaryReader(bytes.NewReader(tt.input))
			got, err := Parse13xBands(&r)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParse137(t *testing.T) {
	input := build13x([]types.Band{{Band: 7, Class: 1}})
	r := readers.NewBinaryReader(bytes.NewReader(input))
	entry, err := Parse137(&r)
	require.NoError(t, err)
	assert.IsType(t, &types.DownlinkEntry{}, entry)
	assert.Equal(t, []types.Band{{Band: 7, Class: 1}}, entry.Bands())
}

func TestParse138(t *testing.T) {
	input := build13x([]types.Band{{Band: 3, Class: 2}})
	r := readers.NewBinaryReader(bytes.NewReader(input))
	entry, err := Parse138(&r)
	require.NoError(t, err)
	assert.IsType(t, &types.UplinkEntry{}, entry)
	assert.Equal(t, []types.Band{{Band: 3, Class: 2}}, entry.Bands())
}

func TestParse20xBands(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []types.Band
		wantErr  bool
	}{
		{
			name:     "single band with mimo",
			input:    build20x([]types.Band{{Band: 7, Class: 1, Mimo: 2}}),
			expected: []types.Band{{Band: 7, Class: 1, Mimo: 2}},
		},
		{
			name: "multiple bands with mimo",
			input: build20x([]types.Band{
				{Band: 7, Class: 1, Mimo: 2},
				{Band: 3, Class: 2, Mimo: 4},
				{Band: 1, Class: 3, Mimo: 1},
			}),
			expected: []types.Band{
				{Band: 7, Class: 1, Mimo: 2},
				{Band: 3, Class: 2, Mimo: 4},
				{Band: 1, Class: 3, Mimo: 1},
			},
		},
		{
			name: "empty slot skipped",
			input: build20x([]types.Band{
				{Band: 0, Class: 0, Mimo: 0},
				{Band: 7, Class: 1, Mimo: 2},
			}),
			expected: []types.Band{{Band: 7, Class: 1, Mimo: 2}},
		},
		{
			name: "invalid class skipped",
			input: build20x([]types.Band{
				{Band: 7, Class: 10, Mimo: 2},
				{Band: 3, Class: 1, Mimo: 4},
			}),
			expected: []types.Band{{Band: 3, Class: 1, Mimo: 4}},
		},
		{
			name:    "truncated input",
			input:   []byte{0x07, 0x00, 0x01, 0x02},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := readers.NewBinaryReader(bytes.NewReader(tt.input))
			got, err := Parse20xBands(&r)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParse201(t *testing.T) {
	input := build20x([]types.Band{{Band: 7, Class: 1, Mimo: 2}})
	r := readers.NewBinaryReader(bytes.NewReader(input))
	entry, err := Parse201(&r)
	require.NoError(t, err)
	assert.IsType(t, &types.DownlinkEntry{}, entry)
	assert.Equal(t, []types.Band{{Band: 7, Class: 1, Mimo: 2}}, entry.Bands())
}

func TestParse202(t *testing.T) {
	input := build20x([]types.Band{{Band: 3, Class: 2, Mimo: 4}})
	r := readers.NewBinaryReader(bytes.NewReader(input))
	entry, err := Parse202(&r)
	require.NoError(t, err)
	assert.IsType(t, &types.UplinkEntry{}, entry)
	assert.Equal(t, []types.Band{{Band: 3, Class: 2, Mimo: 4}}, entry.Bands())
}

func TestParse33xBands(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []types.Band
		wantErr  bool
	}{
		{
			name: "single band with antennas",
			input: build33x([]types.Band{
				{Band: 7, Class: 1, Antennas: []types.Antenna{1, 2, 3}},
			}),
			expected: []types.Band{
				{Band: 7, Class: 1, Antennas: []types.Antenna{1, 2, 3}},
			},
		},
		{
			name: "multiple bands with antennas",
			input: build33x([]types.Band{
				{Band: 7, Class: 1, Antennas: []types.Antenna{1, 2, 3}},
				{Band: 3, Class: 2, Antennas: []types.Antenna{4, 5}},
			}),
			expected: []types.Band{
				{Band: 7, Class: 1, Antennas: []types.Antenna{1, 2, 3}},
				{Band: 3, Class: 2, Antennas: []types.Antenna{4, 5}},
			},
		},
		{
			name: "empty slot skipped",
			input: build33x([]types.Band{
				{Band: 0, Class: 0},
				{Band: 7, Class: 1, Antennas: []types.Antenna{1}},
			}),
			expected: []types.Band{{Band: 7, Class: 1, Antennas: []types.Antenna{1}}},
		},
		{
			name: "invalid class skipped",
			input: build33x([]types.Band{
				{Band: 7, Class: 10, Antennas: []types.Antenna{1}},
				{Band: 3, Class: 1, Antennas: []types.Antenna{2}},
			}),
			expected: []types.Band{{Band: 3, Class: 1, Antennas: []types.Antenna{2}}},
		},
		{
			name:    "truncated input",
			input:   []byte{0x07, 0x00, 0x01, 0x01, 0x02, 0x03},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := readers.NewBinaryReader(bytes.NewReader(tt.input))
			got, err := Parse33xBands(&r)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParse333(t *testing.T) {
	input := build33x([]types.Band{{Band: 7, Class: 1, Antennas: []types.Antenna{1, 2, 3}}})
	r := readers.NewBinaryReader(bytes.NewReader(input))
	entry, err := Parse333(&r)
	require.NoError(t, err)
	assert.IsType(t, &types.DownlinkEntry{}, entry)
	assert.Equal(t, []types.Band{{Band: 7, Class: 1, Antennas: []types.Antenna{1, 2, 3}}}, entry.Bands())
}

func TestParse334(t *testing.T) {
	input := build33x([]types.Band{{Band: 3, Class: 2, Antennas: []types.Antenna{4, 5}}})
	r := readers.NewBinaryReader(bytes.NewReader(input))
	entry, err := Parse334(&r)
	require.NoError(t, err)
	assert.IsType(t, &types.UplinkEntry{}, entry)
	assert.Equal(t, []types.Band{{Band: 3, Class: 2, Antennas: []types.Antenna{4, 5}}}, entry.Bands())
}

func TestParseAntennas(t *testing.T) {
	input := []byte{0x01, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00}
	r := readers.NewBinaryReader(bytes.NewReader(input))
	got := ParseAntennas(&r)
	assert.Equal(t, []types.Antenna{1, 2, 3}, got)
}

func TestParseAntennasIgnoresTrailingZeros(t *testing.T) {
	input := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	r := readers.NewBinaryReader(bytes.NewReader(input))
	got := ParseAntennas(&r)
	assert.Empty(t, got)
}
