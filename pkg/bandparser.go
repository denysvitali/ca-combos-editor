package pkg

import (
	"bufio"
	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseComboText(comboString string) []types.Entry {
	Log.Debugf("comboString: %s", comboString)
	r := readers.NewComboReader(comboString)

	var entries []types.Entry
	dl := types.DownlinkEntry{}
	ul := types.UplinkEntry{}

	cont := true
	for cont {
		band, err := r.ReadNumber()
		if err != nil {
			break
		}

		b := types.Band{}
		b.Band = band
		b.Class = r.ReadClass()

		if b.Class == -1 {
			break
		}

		mimo, err := r.ReadNumber()
		countMimo := len(strconv.Itoa(mimo))

		if err != nil || countMimo == 0 {
			dl.SetBands(append(dl.Bands(), b))
		} else {
			for i := 0; i < countMimo; i++ {
				b.Mimo = mimo
				dl.SetBands(append(dl.Bands(), b))
			}
		}

		ulClass := r.ReadClass()
		if ulClass > 0 {
			mimo, err := r.ReadNumber()
			if err != nil {
				// No MIMO specified (e.g: 41A4A), setting it to 1
				mimo = 1
			}
			ulBand := types.Band{
				Band:     band,
				Class:    ulClass,
				Mimo:     mimo,
			}
			ul.SetBands(append(ul.Bands(), ulBand))
		}

		cont = readers.HasNextBand(&r)
	}

	entries = append(entries, &dl)
	entries = append(entries, &ul)

	if len(ul.Bands()) > 1 {
		Log.Warnf("UL => %v", ul)
		return nil
	}

	Log.Debugf("=> %v", entries)
	return entries
}

func ParseBandDLULFile(downlink string, uplink string) []types.Entry {
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

	var finalEntries []types.Entry

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
		var ulEntries []types.Entry

		sort.Sort(sort.StringSlice(ulBands))

		if len(ulBands) > 0 && ulText != "" {
			for _, bText := range ulBands {
				ulEntry := types.UplinkEntry{}
				ulEntry.SetBands([]types.Band{parseSingleBand(bText)})
				ulEntries = append(ulEntries, &ulEntry)
			}
		}

		finalEntries = append(finalEntries, entry)
		finalEntries = append(finalEntries, ulEntries...)
	}

	return finalEntries
}

func parseSingleBand(text string) types.Band {
	r := readers.NewComboReader(text)

	bandNumber, err := r.ReadNumber()
	if err != nil {
		Log.Fatal(err)
	}

	bandClass := r.ReadClass()
	if bandClass == -1 {
		Log.Fatal("invalid band class")
	}

	return types.Band{Band: bandNumber, Class: bandClass}
}

func ParseBandFile(path string) []types.Entry {
	comboFile, err := os.Open(path)
	if err != nil {
		Log.Fatal(err)
	}
	defer comboFile.Close()

	var finalEntries []types.Entry
	var finalEntriesHM = make(map[string][]types.Entry)

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
		var dlEntries []types.DownlinkEntry
		var ulEntries []types.UplinkEntry

		for _, e := range v {
			switch e := e.(type) {
			case *types.DownlinkEntry:
				dlEntries = append(dlEntries, *e)
			case *types.UplinkEntry:
				ulEntries = append(ulEntries, *e)
			}
		}
		sort.Sort(types.UlArr(ulEntries))

		for _, e := range dlEntries {
			dlEntry := types.DownlinkEntry{}
			dlEntry.SetBands(e.Bands())
			finalEntries = append(finalEntries, &dlEntry)
		}

		for _, e := range ulEntries {
			ulEntry := types.UplinkEntry{}
			ulEntry.SetBands(e.Bands())
			finalEntries = append(finalEntries, &ulEntry)
		}
	}

	if err := scanner.Err(); err != nil {
		Log.Fatal(err)
	}

	return finalEntries
}
