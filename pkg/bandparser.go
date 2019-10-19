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

	classIndex := strings.Index(classes, string(c))

	if classIndex == -1 {
		r.GoBack()
		return -1
	}

	return  classIndex + 1
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

		mimo, err := r.readNumber()
		Log.Debugf("MIMO: %d", mimo)

		countMimo := len(strconv.Itoa(mimo))

		if err != nil || countMimo == 0 {
			dl.bands = append(dl.bands, b)
		} else {
			for i:=0; i< countMimo; i++ {
				dl.bands = append(dl.bands, b)
			}
		}

		ulClass := r.readClass()
		if ulClass > 0 {
			ulband := Band{
				band,
				ulClass,
			}


			if err != nil || countMimo == 0 {
				ul.bands = append(ul.bands, ulband)
			} else {
				for i:=0; i< countMimo; i++ {
					ul.bands = append(ul.bands, ulband)
				}
			}
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
	var finalEntriesHM = make(map[string][]Entry)


	scanner := bufio.NewScanner(comboFile)
	for scanner.Scan() {
		text := scanner.Text()

		if text == "" {
			continue
		}
		entries := parseComboText(text)
		var dl = ""
		for _, e := range entries {
			if dl == "" {
				if e.Name() == "DL" {
					dl = e.String()
				} else {
					continue
				}
			}

			if e.Name() == "DL" && len(finalEntriesHM[dl]) > 1 {
				continue
			}

			finalEntriesHM[dl] = append(finalEntriesHM[dl], e)
		}
	}

	for _, v := range finalEntriesHM {
		finalEntries = append(finalEntries, v...)
	}

	if err := scanner.Err(); err != nil {
		Log.Fatal(err)
	}

	return finalEntries
}