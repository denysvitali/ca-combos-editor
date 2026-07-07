package pkg

import (
	"testing"

	"github.com/denysvitali/ca-combos-editor/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCombo(t *testing.T) {
	comboString := "3A2-1A4A"
	entries, err := parseComboText(comboString)
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "DL", entries[0].Name())
	assert.Equal(t, "UL", entries[1].Name())
}

func TestParseComboMIMO(t *testing.T) {
	comboString := "1A4A-1A4"
	logrus.SetLevel(logrus.DebugLevel)
	entries, err := parseComboText(comboString)
	require.NoError(t, err)

	b1 := types.Band{
		Band:  1,
		Class: 1,
		Mimo:  4,
	}

	b1_2 := types.Band{
		Band:  1,
		Class: 1,
		Mimo:  1,
	}
	assert.Equal(t, []types.Band{b1, b1}, entries[0].Bands())
	assert.Equal(t, []types.Band{b1_2}, entries[1].Bands())
	assert.IsType(t, &types.DownlinkEntry{}, entries[0])
	assert.IsType(t, &types.UplinkEntry{}, entries[1])
}

func TestParseComplexCombo(t *testing.T) {
	Log.Level = logrus.DebugLevel
	entries, err := parseComboText("41A4A-28A2-3A2")
	require.NoError(t, err)

	b1 := types.Band{
		Band:  41,
		Mimo:  4,
		Class: 1,
	}

	b1_2 := types.Band{
		Band:  41,
		Mimo:  1,
		Class: 1,
	}

	b2 := types.Band{
		Band:  28,
		Class: 1,
		Mimo:  2,
	}

	b3 := types.Band{
		Band:  3,
		Class: 1,
		Mimo:  2,
	}

	assert.Equal(t, []types.Band{b1, b2, b3}, entries[0].Bands()) // DL
	assert.Equal(t, []types.Band{b1_2}, entries[1].Bands())       // UL
}

func TestParseCombo2(t *testing.T) {
	Log.Level = logrus.DebugLevel
	entries, err := parseComboText("3C44A-0")
	require.NoError(t, err)

	// DL: 3C 3C
	// UL: 3A
	assert.Equal(t, 2, len(entries))

	dlEntry, ok := entries[0].(*types.DownlinkEntry)
	require.True(t, ok)

	bands := dlEntry.Bands()
	require.Len(t, bands, 2)

	firstBand := bands[0]
	assert.Equal(t, 3, firstBand.Band)
	assert.Equal(t, 3, firstBand.Class) // C

	secondBand := bands[1]
	assert.Equal(t, 3, secondBand.Band)
	assert.Equal(t, 3, secondBand.Class) // C

	ulEntry, ok := entries[1].(*types.UplinkEntry)
	require.True(t, ok)

	thirdBand := ulEntry.Bands()[0]
	assert.Equal(t, 3, thirdBand.Band)
	assert.Equal(t, 1, thirdBand.Class)
}

func TestParseBand1(t *testing.T) {
	comboString := "2A2A-46E2-48C2"
	Log.Level = logrus.DebugLevel
	entries, err := parseComboText(comboString)
	require.NoError(t, err)

	assert.NotNil(t, entries)
	assert.Equal(t, &types.DownlinkEntry{
		BandArr: []types.Band{
			{Band: 48, Class: 3, Mimo: 2},
			{Band: 46, Class: 5, Mimo: 2},
			{Band: 2, Class: 1, Mimo: 2},
		},
	}, entries[0])
	assert.Equal(t, &types.UplinkEntry{
		BandArr: []types.Band{
			{Band: 2, Class: 1, Mimo: 1},
		},
	}, entries[1])
}

func TestParseFile(t *testing.T) {
	Log.Level = logrus.DebugLevel
	entries, err := ParseBandFile("../test/resources/2019-10-17/bands.txt")
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestParseSingleBand(t *testing.T) {
	band, err := parseSingleBand("41A4")
	require.NoError(t, err)
	assert.Equal(t, types.Band{Band: 41, Class: 1, Mimo: 0}, band)
}

func TestParseSingleBandInvalid(t *testing.T) {
	_, err := parseSingleBand("not-a-band")
	require.Error(t, err)
}
