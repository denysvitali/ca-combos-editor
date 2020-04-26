package pkg

import (
	"bufio"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestParseCombo(t *testing.T) {
	comboString := "3A2-1A4A"
	entries := parseComboText(comboString)

	for _, e := range entries {
		log.Printf("Entry: %v", e)
	}
}

func TestParseComboMIMO(t *testing.T) {
	comboString := "1A4A-1A4"
	logrus.SetLevel(logrus.DebugLevel)
	entries := parseComboText(comboString)

	b1 := types.Band {
		Band:  1,
		Class: 1,
		Mimo:  4,
	}

	b1_2 := types.Band {
		Band: 1,
		Class: 1,
		Mimo: 1,
	}
	assert.Equal(t, []types.Band{b1, b1}, entries[0].Bands())
	assert.Equal(t, []types.Band{b1_2}, entries[1].Bands())
	assert.IsType(t, &types.DownlinkEntry{}, entries[0])
	assert.IsType(t, &types.UplinkEntry{}, entries[1])
}

func TestParseComplexCombo(t *testing.T) {
	Log.Level = logrus.DebugLevel
	entries := parseComboText("41A4A-28A2-3A2")

	b1 := types.Band{
		Band: 41,
		Mimo: 4,
		Class: 1,
	}

	b1_2 := types.Band{
		Band: 41,
		Mimo: 1,
		Class: 1,
	}

	b2 := types.Band {
		Band: 28,
		Class: 1,
		Mimo: 2,
	}

	b3 := types.Band{
		Band: 3,
		Class: 1,
		Mimo: 2,
	}

	assert.Equal(t, entries[0].Bands(), []types.Band{b1, b2, b3}) // DL
	assert.Equal(t, entries[1].Bands(), []types.Band{b1_2}) // UL

	log.Printf("Entries: %v", entries)
}

func TestParseCombo2(t *testing.T) {
	Log.Level = logrus.DebugLevel
	entries := parseComboText("3C44A-0")

	// DL: 3C 3C
	// UL: 3A

	assert.Equal(t, 2, len(entries))

	dlEntry, ok := entries[0].(*types.DownlinkEntry)
	assert.True(t, ok)

	bands := dlEntry.Bands()
	assert.Equal(t, 2, len(bands))

	firstBand := bands[0]
	assert.Equal(t, 3, firstBand.Band)
	assert.Equal(t, 3, firstBand.Class) // C

	secondBand := bands[1]
	assert.Equal(t, 3, secondBand.Band)
	assert.Equal(t, 3, secondBand.Class) // C

	ulEntry, ok := entries[1].(*types.UplinkEntry)
	assert.True(t, ok)

	thirdBand := ulEntry.Bands()[0]
	assert.Equal(t, 3, thirdBand.Band)
	assert.Equal(t, 1, thirdBand.Class)

}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func TestParseFile(t *testing.T) {
	Log.Level = logrus.DebugLevel
	ParseBandFile("../test/resources/2019-10-17/bands.txt")
}
