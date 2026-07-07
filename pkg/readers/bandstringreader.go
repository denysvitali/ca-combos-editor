package readers

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// ComboReader parses human-readable combo strings such as "3A2-1A4A".
type ComboReader struct {
	reader *strings.Reader
}

func NewComboReader(s string) ComboReader {
	return ComboReader{reader: strings.NewReader(s)}
}

// NextRune returns the next rune or an error at EOF.
func (r *ComboReader) NextRune() (rune, error) {
	ch, _, err := r.reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return ch, nil
}

// GoBack unreads the last rune.
func (r *ComboReader) GoBack() {
	_ = r.reader.UnreadRune()
}

// ReadNumber reads a consecutive run of decimal digits.
func (r *ComboReader) ReadNumber() (int, error) {
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

// ReadClass reads an LTE bandwidth class letter (A-E) and returns its numeric
// value (1-5). It returns -1 if the next rune is not a recognised class.
func (r *ComboReader) ReadClass() int {
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

	return classIndex + 1
}

// HasNextBand reports whether the next rune is a band separator ('-').
func HasNextBand(r *ComboReader) bool {
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

// Remaining returns the unread portion of the input string.
func (r *ComboReader) Remaining() string {
	rest, _ := io.ReadAll(r.reader)
	return string(rest)
}

// TODO: remove after confirming no callers remain.
func (r *ComboReader) skipOrFailGracefully(expectedRune rune) {}
