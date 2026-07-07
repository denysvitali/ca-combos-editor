package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/denysvitali/ca-combos-editor/pkg/parsers"
	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
	"github.com/sirupsen/logrus"
)

// Log is the package-level logger. Callers may configure its level.
var Log = logrus.New()

// ComboEdit parses an uncompressed NV ITEM 00028874 payload.
type ComboEdit struct {
	FileContent []byte
}

// ComboWriterMode selects the on-wire entry format used when serializing combos.
type ComboWriterMode int

//go:generate go run golang.org/x/tools/cmd/stringer -type=ComboWriterMode
const (
	// COMBOWRITER_137_138 emits 137/138 entries without MIMO bytes.
	COMBOWRITER_137_138 ComboWriterMode = 137
	// COMBOWRITER_201_202 emits 201/202 entries with MIMO bytes.
	COMBOWRITER_201_202 ComboWriterMode = 201
	// COMBOWRITER_333_334 emits 333/334 entries with antenna bytes.
	COMBOWRITER_333_334 ComboWriterMode = 333
)

// ComboWriter serializes a slice of entries to the uncompressed payload format.
type ComboWriter struct {
	FileBody []byte
	Mode     ComboWriterMode
}

// NewComboEdit creates a parser for the given uncompressed payload bytes.
func NewComboEdit(input []byte) ComboEdit {
	return ComboEdit{FileContent: input}
}

// ComboFile is the decoded representation of a payload.
type ComboFile struct {
	EntriesLen uint16
	Entries    []types.Entry
}

// Parse decodes the payload into a ComboFile.
func (c *ComboEdit) Parse() (ComboFile, error) {
	r := readers.NewMyReader(bytes.NewReader(c.FileContent))
	if err := r.Expect(0x00); err != nil {
		return ComboFile{}, fmt.Errorf("header byte 0: %w", err)
	}
	if err := r.Expect(0x00); err != nil {
		return ComboFile{}, fmt.Errorf("header byte 1: %w", err)
	}

	cf := ComboFile{}

	lenArr, err := r.ReadBytes(2)
	if err != nil {
		return ComboFile{}, fmt.Errorf("header length: %w", err)
	}
	cf.EntriesLen = binary.LittleEndian.Uint16(lenArr)
	Log.Info("This CA bands file contains ", cf.EntriesLen, " entries")

	for i := uint16(0); i < cf.EntriesLen; i++ {
		entry, err := c.parseEntry(&r)
		if err != nil {
			return ComboFile{}, fmt.Errorf("entry %d: %w", i, err)
		}
		cf.Entries = append(cf.Entries, entry)
	}

	return cf, nil
}

// Write serializes the provided entries and returns the full payload bytes.
func (w *ComboWriter) Write(entries []types.Entry) []byte {
	output := make([]byte, 0, types.HeaderSize+len(entries)*(2+types.MaxBandsPerEntry*w.bandRecordSize()))
	output = append(output, 0x00, 0x00)

	countBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(countBytes, uint16(len(entries)))
	output = append(output, countBytes...)

	for _, e := range entries {
		w.writeEntry(e)
	}

	output = append(output, w.FileBody...)

	Log.Infof("Writing %d entries...", len(entries))
	return output
}

// ReadComboFile parses an uncompressed combo file and prints its entries.
func ReadComboFile(path string) error {
	result, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read combo file: %w", err)
	}

	ce := NewComboEdit(result)
	cf, err := ce.Parse()
	if err != nil {
		return fmt.Errorf("parse combo file: %w", err)
	}

	for _, e := range cf.Entries {
		fmt.Printf("%s: %v\n", e.Name(), e)
	}
	return nil
}

// WriteComboFile serializes entries to path using the selected mode.
func WriteComboFile(entries []types.Entry, mode ComboWriterMode, path string) (err error) {
	w := ComboWriter{Mode: mode}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to open file %q for writing: %w", path, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close file %q: %w", path, cerr)
		}
	}()

	if _, err := io.Copy(f, bytes.NewReader(w.Write(entries))); err != nil {
		return fmt.Errorf("unable to write file %q: %w", path, err)
	}
	return nil
}

func (c *ComboEdit) parseEntry(r *readers.BinaryReader) (types.Entry, error) {
	typeBytes, err := r.ReadBytes(2)
	if err != nil {
		return nil, fmt.Errorf("read entry type: %w", err)
	}
	entryType := types.EntryType(binary.LittleEndian.Uint16(typeBytes))
	Log.Debugf("Parsing entry type %d", entryType)

	switch entryType {
	case types.EntryTypeDownlinkNoMIMO:
		return parsers.Parse137(r)
	case types.EntryTypeUplinkNoMIMO:
		return parsers.Parse138(r)
	case types.EntryTypeDownlinkMIMO:
		return parsers.Parse201(r)
	case types.EntryTypeUplinkMIMO:
		return parsers.Parse202(r)
	case types.EntryTypeDownlinkAntennas:
		return parsers.Parse333(r)
	case types.EntryTypeUplinkAntennas:
		return parsers.Parse334(r)
	default:
		return nil, fmt.Errorf("invalid entry type %d", entryType)
	}
}

func (w *ComboWriter) writeEntry(entry types.Entry) {
	switch entry := entry.(type) {
	case *types.DownlinkEntry:
		w.writeType(w.downlinkType())
		sortedBands := make([]types.Band, len(entry.Bands()))
		copy(sortedBands, entry.Bands())
		sort.Sort(sort.Reverse(types.BandArr(sortedBands)))
		w.writeBands(sortedBands)

	case *types.UplinkEntry:
		w.writeType(w.uplinkType())
		w.writeBands(entry.Bands())
	}
}

func (w *ComboWriter) writeType(t types.EntryType) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(t))
	w.FileBody = append(w.FileBody, b...)
}

func (w *ComboWriter) downlinkType() types.EntryType {
	switch w.Mode {
	case COMBOWRITER_137_138:
		return types.EntryTypeDownlinkNoMIMO
	case COMBOWRITER_201_202:
		return types.EntryTypeDownlinkMIMO
	case COMBOWRITER_333_334:
		return types.EntryTypeDownlinkAntennas
	}
	return types.EntryTypeDownlinkNoMIMO
}

func (w *ComboWriter) uplinkType() types.EntryType {
	switch w.Mode {
	case COMBOWRITER_137_138:
		return types.EntryTypeUplinkNoMIMO
	case COMBOWRITER_201_202:
		return types.EntryTypeUplinkMIMO
	case COMBOWRITER_333_334:
		return types.EntryTypeUplinkAntennas
	}
	return types.EntryTypeUplinkNoMIMO
}

func (w *ComboWriter) writeBands(bands []types.Band) {
	for i := range types.MaxBandsPerEntry {
		if i < len(bands) {
			w.writeBand(bands[i])
		} else {
			w.writeEmptyBand()
		}
	}
}

func (w *ComboWriter) writeBand(b types.Band) {
	bandBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bandBytes, uint16(b.Band))
	w.FileBody = append(w.FileBody, bandBytes...)
	w.FileBody = append(w.FileBody, byte(b.Class))

	switch w.Mode {
	case COMBOWRITER_201_202:
		w.FileBody = append(w.FileBody, byte(b.Mimo))
	case COMBOWRITER_333_334:
		antennaBytes := make([]byte, types.AntennaCount)
		for i, a := range b.Antennas {
			if i >= types.AntennaCount {
				break
			}
			antennaBytes[i] = byte(a)
		}
		w.FileBody = append(w.FileBody, antennaBytes...)
	}
}

func (w *ComboWriter) writeEmptyBand() {
	empty := make([]byte, w.bandRecordSize())
	w.FileBody = append(w.FileBody, empty...)
}

func (w *ComboWriter) bandRecordSize() int {
	switch w.Mode {
	case COMBOWRITER_201_202:
		return types.BandRecordSize20x
	case COMBOWRITER_333_334:
		return types.BandRecordSize33x
	}
	return types.BandRecordSize13x
}

// SetMode configures the writer mode.
func (w *ComboWriter) SetMode(mode ComboWriterMode) {
	w.Mode = mode
}
