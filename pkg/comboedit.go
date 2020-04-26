package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/denysvitali/ca-combos-editor/pkg/parsers"
	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
)

var Log = logrus.New()

type ComboEdit struct {
	FileContent []byte
}

type ComboWriterMode int

//go:generate stringer -type=ComboWriterMode
const (
	COMBOWRITER_137_138 ComboWriterMode = 137
	COMBOWRITER_201_202 ComboWriterMode = 201
)

type ComboWriter struct {
	FileBody      []byte
	EntriesLength int
	Mode          ComboWriterMode
}

func NewComboEdit(input []byte) ComboEdit {
	return ComboEdit{FileContent: input}
}

type ComboFile struct {
	EntriesLen uint16
	Entries    []types.Entry
}

func (c *ComboEdit) Parse() ComboFile {
	r := readers.NewMyReader(bytes.NewReader(c.FileContent))
	r.Expect(0x00)
	r.Expect(0x00)

	cf := ComboFile{}

	// Length Section
	var lenArr []byte
	lenArr = append(lenArr, r.Rb())
	lenArr = append(lenArr, r.Rb())
	cf.EntriesLen = binary.LittleEndian.Uint16(lenArr)
	Log.Info("This CA Bands files contains ", cf.EntriesLen, " entries")

	for i := uint16(0); i < cf.EntriesLen; i++ {
		cf.Entries = append(cf.Entries, c.parseEntry(&r))
	}

	return cf
}

func (w *ComboWriter) Write(entries []types.Entry) []byte {
	var output []byte
	output = append(output, 0x00)
	output = append(output, 0x00)

	var b = make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(len(entries)))
	output = append(output, b...)

	for _, e := range entries {
		w.writeEntry(e)
	}

	output = append(output, w.FileBody...)

	Log.Infof("Writing %d entries...", len(entries))

	return output
}

func ReadComboFile(path string) {
	result, err := ioutil.ReadFile(path)

	if err != nil {
		Log.Fatal(err)
	}

	ce := NewComboEdit(result)
	cf := ce.Parse()

	for _, e := range cf.Entries {
		fmt.Printf("%s: %v\n", e.Name(), e)
	}
}


func WriteComboFile(entries []types.Entry, mode ComboWriterMode, path string) {
	w := ComboWriter{}
	w.SetMode(mode)

	f, err := os.Create(path)
	if err != nil {
		Log.Fatalf("unable to open file \"%s\" for writing", path)
		return
	}
	defer f.Close()
	_, err = f.Write(w.Write(entries))
	if err != nil {
		Log.Fatalf("unable to write file \"%s\"", path)
	}
}

func (c *ComboEdit) parseEntry(r *readers.BinaryReader) types.Entry {
	var e types.Entry
	entryType := binary.LittleEndian.Uint16(r.ReadBytes(2))
	Log.Debugf("Parsing entry type %d", entryType)

	switch entryType {
	case 137:
		e = parsers.Parse137(r)
	case 138:
		e = parsers.Parse138(r)
	case 201:
		e = parsers.Parse201(r)
	case 202:
		e = parsers.Parse202(r)
	case 333:
		e = parsers.Parse333(r)
	case 334:
		e = parsers.Parse334(r)
	default:
		Log.Errorf("Invalid type %d found!", entryType)
	}
	return e
}

func (w *ComboWriter) writeEntry(entry types.Entry) {
	switch entry.(type) {
	case *types.DownlinkEntry:
		switch w.Mode {
		case COMBOWRITER_137_138:
			w.FileBody = append(w.FileBody, byte(137))
		case COMBOWRITER_201_202:
			w.FileBody = append(w.FileBody, byte(201))
		}
		w.FileBody = append(w.FileBody, 0)
		sortedBands := entry.Bands()
		sort.Sort(types.BandArr(sortedBands))

		for i := len(sortedBands)/2 - 1; i >= 0; i-- {
			opp := len(sortedBands) - 1 - i
			sortedBands[i], sortedBands[opp] = sortedBands[opp], sortedBands[i]
		}

		w.writeBands(sortedBands)

	case *types.UplinkEntry:
		switch w.Mode {
		case COMBOWRITER_137_138:
			w.FileBody = append(w.FileBody, byte(138))
		case COMBOWRITER_201_202:
			w.FileBody = append(w.FileBody, byte(202))
		}
		w.FileBody = append(w.FileBody, 0)
		w.writeBands(entry.Bands())
	}
}

func (w *ComboWriter) writeBands(bands []types.Band) {
	for i := 0; i < 6; i++ {
		if i < len(bands) {
			w.FileBody = append(w.FileBody, byte(bands[i].Band))
			w.FileBody = append(w.FileBody, 0x00)
			w.FileBody = append(w.FileBody, byte(bands[i].Class))

			if w.Mode == COMBOWRITER_201_202 {
				w.FileBody = append(w.FileBody, byte(bands[i].Mimo))
			}
		} else {
			w.FileBody = append(w.FileBody, 0x00)
			w.FileBody = append(w.FileBody, 0x00)
			w.FileBody = append(w.FileBody, 0x00)
			if w.Mode == COMBOWRITER_201_202 {
				w.FileBody = append(w.FileBody, 0x00)
			}
		}
	}
}

func (w *ComboWriter) SetMode(mode ComboWriterMode) {
	w.Mode = mode
}
