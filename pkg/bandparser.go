package pkg

import (
	"bufio"
	"errors"
	"github.com/Sirupsen/logrus"
	"os"
	"sort"
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
	Log.Debugf("comboString: %s", comboString)
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

		if b.Class == -1 {
			break
		}

		mimo, err := r.readNumber()
		//Log.Debugf("MIMO: %d", mimo)

		countMimo := len(strconv.Itoa(mimo))

		if err != nil || countMimo == 0 {
			dl.bands = append(dl.bands, b)
		} else {
			for i:=0; i < countMimo; i++ {
					dl.bands = append(dl.bands, b)
			}
		}

		ulClass := r.readClass()
		if ulClass > 0 {
			ulBand := Band{
				band,
				ulClass,
				0,
			}
			ul.bands = append(ul.bands, ulBand)
		}

		cont = hasNextBand(&r)
	}

	entries = append(entries, &dl)
	entries = append(entries, &ul)

	if len(ul.bands) > 1 {
		Log.Warnf("UL => %v", ul)
		return nil
	}

	Log.Debugf("=> %v", entries)
	return entries
}

func ParseBandDLULFile(downlink string, uplink string) []Entry {
	dlFile, err := os.Open(downlink)
	if err != nil {
		Log.Fatal(err)
	}
	ulFile, err := os.Open(uplink)
	if err != nil {
		Log.Fatal(err)
	}

	defer dlFile.Close()
	defer ulFile.Close()

	var finalEntries []Entry

	dlScanner := bufio.NewScanner(dlFile)
	ulScanner := bufio.NewScanner(ulFile)

	for dlScanner.Scan() {
		ulScanner.Scan()
		dlText := dlScanner.Text()
		ulText := ulScanner.Text()

		if dlText == "" {
			continue
		}

		entry := parseComboText(dlText)[0]
		ulBands := strings.Split(ulText, ", ")
		var ulEntries []Entry

		sort.Sort(sort.StringSlice(ulBands))


		if len(ulBands) > 0 && ulText != "" {
			for _, bText := range ulBands {
				ulEntries = append(ulEntries, &UplinkEntry{bands: []Band{parseSingleBand(bText)}})
			}
		}

		finalEntries = append(finalEntries, entry)
		finalEntries = append(finalEntries, ulEntries...)
	}

	return finalEntries
}

func parseSingleBand(text string) Band {
	r := NewMyStringReader(text)

	bandNumber, err := r.readNumber()
	if err != nil {
		logrus.Fatal(err)
	}

	bandClass := r.readClass()
	if bandClass == -1 {
		logrus.Fatal("invalid band class")
	}

	return Band{Band:bandNumber, Class: bandClass}
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

			found := false
			for _, v := range finalEntriesHM[dl] {
				if v.String() == e.String() {
					found = true
				}
			}

			if found {
				continue
			}

			finalEntriesHM[dl] = append(finalEntriesHM[dl], e)
		}
	}

	for _, v := range finalEntriesHM {
		var dlEntries []DownlinkEntry
		var ulEntries []UplinkEntry

		for _, e := range v {
			switch e := e.(type) {
			case *DownlinkEntry:
				dlEntries = append(dlEntries, *e)
			case *UplinkEntry:
				ulEntries = append(ulEntries, *e)
			}
		}
		sort.Sort(UlArr(ulEntries))

		for _, e := range dlEntries {
			finalEntries = append(finalEntries, &DownlinkEntry{bands: e.bands})
		}

		for _, e := range ulEntries {
			finalEntries = append(finalEntries, &UplinkEntry{bands: e.bands})
		}



	}

	if err := scanner.Err(); err != nil {
		Log.Fatal(err)
	}

	return finalEntries
}