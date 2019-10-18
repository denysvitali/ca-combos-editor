package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
)

var Log = logrus.New()

type ComboEdit struct {
	FileContent []byte
}

type ComboWriter struct {
	FileBody []byte
	EntriesLength int
}

func NewComboEdit(input []byte) ComboEdit {
	return ComboEdit{FileContent:input}
}

type MyReader struct {
	reader *bytes.Reader
}

func NewMyReader(reader *bytes.Reader) MyReader {
	return MyReader{reader:reader}
}

func (m *MyReader) rb() byte {
	b, e := m.reader.ReadByte()
	if e != nil {
		Log.Fatal(e)
	}

	return b
}

func (m *MyReader) expect(b byte) {
	found := m.rb()
	if found != b {
		Log.Fatalf("Unexpected byte %02X found, %02X expected",
				found & 0xFF,
				b & 0xFF)
	}
}

type ComboFile struct {
	Entries_Len uint16
	Entries     []Entry
}

type Entry interface {
	Name() string
	Bands() []Band
}

type UplinkEntry struct {
	bands []Band
}

func (u* UplinkEntry) Bands() []Band {
	return u.bands
}

func (u* UplinkEntry) Name() string {
	return "UL"
}

type DownlinkEntry struct {
	bands []Band
}

func (d* DownlinkEntry) Bands() []Band {
	return d.bands;
}

func (d* DownlinkEntry) Name() string {
	return "DL";
}

type Band struct {
	Band  int
	Class int
}

func (c *ComboEdit) Parse() ComboFile {
	r := NewMyReader(bytes.NewReader(c.FileContent))
	r.expect(0x00)
	r.expect(0x00)

	cf := ComboFile{}

	// Length Section
	var lenArr []byte
	lenArr = append(lenArr, r.rb())
	lenArr = append(lenArr, r.rb())
	cf.Entries_Len = binary.LittleEndian.Uint16(lenArr)
	Log.Info("This CA Bands files contains ", cf.Entries_Len, " entries")

	for i:=uint16(0); i<cf.Entries_Len; i++ {
		cf.Entries = append(cf.Entries, c.parseEntry(&r))
	}

	return cf
}

func (w *ComboWriter) Write(entries []Entry) []byte {
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
		fmt.Printf("Entry %s: %v\n", e.Name(), e.Bands())
	}
}

func WriteComboFile(entries []Entry, path string) {
	w := ComboWriter{}

	f, err := os.Create(path)
	if err != nil {
		Log.Fatalf("unable to open file \"%s\" for writing", path)
	}
	defer f.Close()
	_, err = f.Write(w.Write(entries))
	if err != nil {
		Log.Fatalf("unable to write file \"%s\"", path)
	}
}


func (c *ComboEdit) parseEntry(r *MyReader) Entry {
	var e Entry
	entryType := int(r.rb())
	r.rb()
	Log.Debugf("Parsing entry type %d", entryType)

	switch entryType {
	case 137:
		// DL
		dlEntry := DownlinkEntry{}
		dlEntry.bands =  parseBands(r)
		e = &dlEntry
	case 138:
		// UL
		ulEntry := UplinkEntry{}
		ulEntry.bands = parseBands(r)
		e = &ulEntry
	default:
		Log.Warnf("Invalid type %d found!", entryType)
	}
	return e
}

func (c *ComboWriter) writeEntry(entry Entry){
	switch entry.(type) {
	case *DownlinkEntry:
		c.FileBody = append(c.FileBody, byte(137))
		c.FileBody = append(c.FileBody, 0)
		c.writeBands(entry.Bands())


	case *UplinkEntry:
		c.FileBody = append(c.FileBody, byte(138))
		c.FileBody = append(c.FileBody, 0)
		c.writeBands(entry.Bands())
	}
}

func parseBands(r *MyReader) []Band {
	var combos []Band

	for i:=0; i<6;i++ {
		bwc := Band{}
		band := int(r.rb())
		r.expect(0x00)
		class := int(r.rb())

		bwc.Band = band
		bwc.Class = class

		combos = append(combos, bwc)
	}
	return combos
}

func (c *ComboWriter) writeBands(bands []Band){
	for i := 0; i<6; i++ {
		if i < len(bands) {
			c.FileBody = append(c.FileBody, byte(int8(bands[i].Band)))
			c.FileBody = append(c.FileBody, 0x00)
			c.FileBody = append(c.FileBody, byte(int8(bands[i].Class)))
		} else {
			c.FileBody = append(c.FileBody, 0x00)
			c.FileBody = append(c.FileBody, 0x00)
			c.FileBody = append(c.FileBody, 0x00)
		}
	}
}