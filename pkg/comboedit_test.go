package pkg

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/denysvitali/ca-combos-editor/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	Log.SetLevel(logrus.ErrorLevel)
	os.Exit(m.Run())
}

func TestComboEditParse(t *testing.T) {
	tests := []struct {
		name        string
		fixture     string
		wantEntries int
	}{
		{
			name:        "13x fixture",
			fixture:     "../test/resources/2019-10-17/extracted",
			wantEntries: 193,
		},
		{
			name:        "20x fixture 1",
			fixture:     "../test/resources/2019-11-26/extracted/1",
			wantEntries: 2184,
		},
		{
			name:        "20x fixture 2",
			fixture:     "../test/resources/2019-11-26/extracted/2",
			wantEntries: 445,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.fixture)
			require.NoError(t, err)

			ce := NewComboEdit(data)
			cf, err := ce.Parse()
			require.NoError(t, err)

			assert.Equal(t, tt.wantEntries, len(cf.Entries))
		})
	}
}

func TestComboEditRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		mode    ComboWriterMode
	}{
		{
			name:    "13x fixture round-trip",
			fixture: "../test/resources/2019-10-17/extracted",
			mode:    COMBOWRITER_137_138,
		},
		{
			name:    "20x fixture 1 round-trip",
			fixture: "../test/resources/2019-11-26/extracted/1",
			mode:    COMBOWRITER_201_202,
		},
		{
			name:    "20x fixture 2 round-trip",
			fixture: "../test/resources/2019-11-26/extracted/2",
			mode:    COMBOWRITER_201_202,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.fixture)
			require.NoError(t, err)

			ce := NewComboEdit(data)
			cf, err := ce.Parse()
			require.NoError(t, err)

			w := ComboWriter{Mode: tt.mode}
			serialized := w.Write(cf.Entries)

			ce2 := NewComboEdit(serialized)
			cf2, err := ce2.Parse()
			require.NoError(t, err)

			assert.Equal(t, normalizeEntries(cf.Entries), normalizeEntries(cf2.Entries))
		})
	}
}

func TestComboEditRoundTripAntennaMode(t *testing.T) {
	cf := ComboFile{
		EntriesLen: 2,
		Entries: []types.Entry{
			&types.DownlinkEntry{
				BandArr: []types.Band{
					{Band: 3, Class: 1, Antennas: []types.Antenna{1, 2}},
					{Band: 7, Class: 2, Antennas: []types.Antenna{1, 2, 4}},
				},
			},
			&types.UplinkEntry{
				BandArr: []types.Band{
					{Band: 1, Class: 1, Antennas: []types.Antenna{1}},
					{Band: 20, Class: 1, Antennas: []types.Antenna{2, 4}},
				},
			},
		},
	}

	w := ComboWriter{Mode: COMBOWRITER_333_334}
	serialized := w.Write(cf.Entries)

	ce := NewComboEdit(serialized)
	cf2, err := ce.Parse()
	require.NoError(t, err)

	assert.Equal(t, len(cf.Entries), len(cf2.Entries))
	assert.Equal(t, normalizeEntries(cf.Entries), normalizeEntries(cf2.Entries))
}

func TestComboEditParseInvalidHeader(t *testing.T) {
	ce := NewComboEdit([]byte{0x01, 0x00, 0x00, 0x00})
	_, err := ce.Parse()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "header byte 0")
}

func TestComboEditParseInvalidEntryType(t *testing.T) {
	ce := NewComboEdit([]byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00})
	_, err := ce.Parse()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid entry type")
}

func normalizeEntries(entries []types.Entry) []types.Entry {
	normalized := make([]types.Entry, len(entries))
	for i, e := range entries {
		switch e := e.(type) {
		case *types.DownlinkEntry:
			bands := make([]types.Band, len(e.Bands()))
			copy(bands, e.Bands())
			types.SortBandsAsc(bands)
			normalized[i] = &types.DownlinkEntry{BandArr: bands}
		case *types.UplinkEntry:
			bands := make([]types.Band, len(e.Bands()))
			copy(bands, e.Bands())
			normalized[i] = &types.UplinkEntry{BandArr: bands}
		default:
			normalized[i] = e
		}
	}
	return normalized
}

func TestReadComboFile(t *testing.T) {
	var buf bytes.Buffer
	err := ReadComboFile("../test/resources/2019-10-17/extracted", &buf)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "DL")
	assert.Contains(t, out, "UL")
}

func TestWriteComboFile(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.bin")
	entries, err := ParseBandFile("../test/resources/2019-10-17/bands.txt")
	require.NoError(t, err)

	require.NoError(t, WriteComboFile(entries, COMBOWRITER_137_138, out))
	assert.FileExists(t, out)

	written, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.NotEmpty(t, written)
}

func TestSetMode(t *testing.T) {
	w := ComboWriter{}
	w.SetMode(COMBOWRITER_201_202)
	assert.Equal(t, COMBOWRITER_201_202, w.Mode)
}

func TestComboEditRoundTrip333(t *testing.T) {
	entries := []types.Entry{
		&types.DownlinkEntry{BandArr: []types.Band{
			{Band: 7, Class: 1, Antennas: []types.Antenna{1, 2, 3}},
		}},
		&types.UplinkEntry{BandArr: []types.Band{
			{Band: 7, Class: 1, Antennas: []types.Antenna{1}},
		}},
	}

	w := ComboWriter{Mode: COMBOWRITER_333_334}
	serialized := w.Write(entries)

	ce := NewComboEdit(serialized)
	cf, err := ce.Parse()
	require.NoError(t, err)
	require.Len(t, cf.Entries, 2)

	dl, ok := cf.Entries[0].(*types.DownlinkEntry)
	require.True(t, ok)
	assert.Equal(t, []types.Antenna{1, 2, 3}, dl.Bands()[0].Antennas)

	ul, ok := cf.Entries[1].(*types.UplinkEntry)
	require.True(t, ok)
	assert.Equal(t, []types.Antenna{1}, ul.Bands()[0].Antennas)
}
