package pkg

import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

type MyStringReader struct {
	reader *strings.Reader
}

func NewMyStringReader(string string) MyStringReader{
	return MyStringReader{reader: strings.NewReader(string)}
}

func (r* MyStringReader) NextRune() rune {
	ch, _, error := r.reader.ReadRune()
	if error != nil {
		Log.Fatalf("Unable to get next rune: %s", error)
	}

	return ch
}

func (r* MyStringReader) GoBack() {
	_ = r.reader.UnreadRune()
}

func (r* MyStringReader) readNumber() (int, error) {
	var numberRunes []rune
	c := r.NextRune()
	for unicode.IsNumber(c) {
		numberRunes = append(numberRunes, c)
		c = r.NextRune()
	}

	r.GoBack()

	if len(numberRunes) > 0 {
		number, err := strconv.Atoi(string(numberRunes))
		if err != nil {
			return -1, err
		}

		return number, nil
	}

	return -1, errors.New("number not found")
}

func (r* MyStringReader) readClass() int {
	c := r.NextRune()
	classes := "ABCD"

	return strings.Index(classes, string(c)) + 1
}


func parseComboText(comboString string) Entry {
	r := NewMyStringReader(comboString)

	band, err := r.readNumber()
	if err != nil {
		Log.Fatalf("Unable to parse Combo Text: %s", err)
	}

	b := Band{}
	b.Band = band
	b.Class = r.readClass()

	e := UplinkEntry{}

	return &e
}

func TestParseCombo(t *testing.T){
	comboString := "3A2-1A4A"
	parseComboText(comboString)
}