package pkg

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
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

func (r* MyStringReader) skipOrFailGracefully(expectedRune rune) {

}

func (r* MyStringReader) readClass() int {
	c, _, err := r.reader.ReadRune()
	if err != nil {
		return -1
	}
	classes := "ABCD"

	class_index := strings.Index(classes, string(c))

	if class_index == -1 {
		r.GoBack()
		return -1
	}

	return  class_index + 1
}

func hasNextBand(r* MyStringReader) bool {
	ch, _, err := r.reader.ReadRune()
	if err != nil {
		return false
	}

	if ch == rune('-') {
		return true
	}

	r.reader.UnreadRune()
	return false
}

func parseComboText(comboString string) []Entry {
	r := NewMyStringReader(comboString)

	var entries []Entry
	dl := DownlinkEntry{}
	ul := UplinkEntry{}

	cont := true
	for cont {
		band, err := r.readNumber()
		if err != nil {
			break
		}

		b := Band{}
		b.Band = band
		b.Class = r.readClass()

		dl.bands = append(dl.bands, b)

		mimo, err := r.readNumber()
		Log.Debugf("MIMO: %d", mimo)

		ulClass := r.readClass()
		if ulClass > 0 {
			ulband := Band{
				band,
				ulClass,
			}
			ul.bands = append(ul.bands, ulband)
		}

		cont = hasNextBand(&r)
	}

	entries = append(entries, &dl)
	entries = append(entries, &ul)

	return entries
}

func ParseBandFile(path string) []Entry {
	comboFile, err := os.Open(path)
	if err != nil {
		Log.Fatal(err)
	}
	defer comboFile.Close()

	var finalEntries []Entry

	scanner := bufio.NewScanner(comboFile)
	for scanner.Scan() {
		text := scanner.Text()

		if text == "" {
			continue
		}
		entries := parseComboText(text)
		for _, e := range entries {
			finalEntries = append(finalEntries, e)
		}
	}

	if err := scanner.Err(); err != nil {
		Log.Fatal(err)
	}

	return finalEntries
}