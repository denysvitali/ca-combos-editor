package pkg

import (
	"bufio"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestParseCombo(t *testing.T){
	comboString := "3A2-1A4A"
	entries := parseComboText(comboString)

	for _, e := range entries {
		log.Printf("Entry: %v", e)
	}
}

func TestParseComplexCombo(t *testing.T){
	Log.Level = logrus.DebugLevel
	entries := parseComboText("7C44C-3C22")

	log.Printf("Entries: %v", entries)
}

func TestParseCombo2(t *testing.T) {
	Log.Level = logrus.DebugLevel
	entries := parseComboText("3C44A-0")

	// DL: 3C 3C
	// UL: 3A

	assert.Equal(t,2, len(entries))

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
	var (isPrefix bool = true
		err error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln),err
}

func TestParseFile(t *testing.T){
	Log.Level = logrus.DebugLevel
	ParseBandFile("../test/resources/2019-10-17/bands.txt")
}