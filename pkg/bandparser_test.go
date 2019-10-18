package pkg

import (
	"bufio"
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