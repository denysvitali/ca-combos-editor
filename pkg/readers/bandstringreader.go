package readers

import (
	"errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"unicode"
)

type ComboReader struct {
	reader *strings.Reader
}

func NewComboReader(string string) ComboReader {
	return ComboReader{reader: strings.NewReader(string)}
}

func (r*ComboReader) NextRune() rune {
	ch, _, err := r.reader.ReadRune()
	if err != nil {
		logrus.Fatalf("Unable to get next rune: %s", err)
	}

	return ch
}

func (r*ComboReader) GoBack() {
	_ = r.reader.UnreadRune()
}

func (r*ComboReader) ReadNumber() (int, error) {
	var numberRunes []rune
	c, _, err := r.reader.ReadRune()
	if err != nil {
		return -1, err
	}
	for unicode.IsNumber(c) {
		numberRunes = append(numberRunes, c)
		c, _, err = r.reader.ReadRune()
		if err != nil {
			break
		}
	}

	if err == nil {
		r.GoBack()
	}

	if len(numberRunes) > 0 {
		number, err := strconv.Atoi(string(numberRunes))
		if err != nil {
			return -1, err
		}

		return number, nil
	}

	return -1, errors.New("number not found")
}

func (r*ComboReader) skipOrFailGracefully(expectedRune rune) {

}

func (r*ComboReader) ReadClass() int {
	c, _, err := r.reader.ReadRune()
	if err != nil {
		return -1
	}
	classes := "ABCDE"

	classIndex := strings.Index(classes, string(c))

	if classIndex == -1 {
		r.GoBack()
		return -1
	}

	return  classIndex + 1
}

func HasNextBand(r*ComboReader) bool {
	ch, _, err := r.reader.ReadRune()
	if err != nil {
		return false
	}
	if ch == rune('-') {
		return true
	}

	_ = r.reader.UnreadRune()
	return false
}

