package pkg

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
)

// parseComboText parses a single human-readable combo line into a downlink and
// an optional uplink entry.
func parseComboText(comboString string) ([]types.Entry, error) {
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
			for range countMimo {
				b.Mimo = mimo
				dl.SetBands(append(dl.Bands(), b))
			}
		}

		ulClass := r.ReadClass()
		if ulClass > 0 {
			mimo, err := r.ReadNumber()
			if err != nil {
				// No MIMO specified (e.g: 41A4A), default to 1.
				mimo = 1
			}
			ulBand := types.Band{
				Band:  band,
				Class: ulClass,
				Mimo:  mimo,
			}
			ul.SetBands(append(ul.Bands(), ulBand))
		}

		cont = readers.HasNextBand(&r)
	}

	entries = append(entries, &dl)
	entries = append(entries, &ul)

	if len(ul.Bands()) > 1 {
		return nil, nil
	}

	Log.Debugf("= %v", entries)
	return entries, nil
}

// ParseBandDLULFile creates combo entries from separate downlink and uplink
// description files.
func ParseBandDLULFile(downlink string, uplink string) (finalEntries []types.Entry, err error) {
	dlFile, err := os.Open(downlink)
	if err != nil {
		return nil, fmt.Errorf("open downlink file: %w", err)
	}
	defer func() {
		if cerr := dlFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close downlink file: %w", cerr)
		}
	}()

	ulFile, err := os.Open(uplink)
	if err != nil {
		return nil, fmt.Errorf("open uplink file: %w", err)
	}
	defer func() {
		if cerr := ulFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close uplink file: %w", cerr)
		}
	}()

	dlScanner := bufio.NewScanner(dlFile)
	ulScanner := bufio.NewScanner(ulFile)

	lineNo := 0
	for dlScanner.Scan() {
		lineNo++
		if !ulScanner.Scan() {
			return nil, fmt.Errorf("downlink file has more lines than uplink file at line %d", lineNo)
		}
		dlText := dlScanner.Text()
		ulText := ulScanner.Text()

		if dlText == "" {
			continue
		}

		entries, err := parseComboText(dlText)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
		if len(entries) == 0 {
			continue
		}

		entry := entries[0]
		ulBands := strings.Split(ulText, ", ")
		var ulEntries []types.Entry

		sort.Strings(ulBands)

		if len(ulBands) > 0 && ulBands[0] != "" {
			for _, bText := range ulBands {
				band, err := parseSingleBand(bText)
				if err != nil {
					return nil, fmt.Errorf("line %d uplink %q: %w", lineNo, bText, err)
				}
				ulEntry := types.UplinkEntry{}
				ulEntry.SetBands([]types.Band{band})
				ulEntries = append(ulEntries, &ulEntry)
			}
		}

		finalEntries = append(finalEntries, entry)
		finalEntries = append(finalEntries, ulEntries...)
	}

	if err := dlScanner.Err(); err != nil {
		return nil, fmt.Errorf("read downlink file: %w", err)
	}
	if err := ulScanner.Err(); err != nil {
		return nil, fmt.Errorf("read uplink file: %w", err)
	}

	return finalEntries, nil
}

func parseSingleBand(text string) (types.Band, error) {
	r := readers.NewComboReader(text)

	bandNumber, err := r.ReadNumber()
	if err != nil {
		return types.Band{}, fmt.Errorf("read band number in %q: %w", text, err)
	}

	bandClass := r.ReadClass()
	if bandClass == -1 {
		return types.Band{}, fmt.Errorf("invalid band class in %q", text)
	}

	return types.Band{Band: bandNumber, Class: bandClass}, nil
}

// ParseBandFile parses a bands.txt-style file into a sorted list of entries.
func ParseBandFile(path string) (finalEntries []types.Entry, err error) {
	comboFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open band file: %w", err)
	}
	defer func() {
		if cerr := comboFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close band file: %w", cerr)
		}
	}()

	finalEntriesHM := make(map[string][]types.Entry)

	scanner := bufio.NewScanner(comboFile)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		text := scanner.Text()

		if text == "" {
			continue
		}
		entries, err := parseComboText(text)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}

		dl := ""
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
					break
				}
			}

			if found {
				continue
			}

			finalEntriesHM[dl] = append(finalEntriesHM[dl], e)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read band file: %w", err)
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

	return finalEntries, nil
}
