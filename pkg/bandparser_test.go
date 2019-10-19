package pkg

import (
	"bufio"
	"github.com/Sirupsen/logrus"
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
	ParseBandFile("../test/resources/2019-10-17/bands.txt")
}