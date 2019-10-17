package pkg

import (
	"bytes"
	"github.com/Sirupsen/logrus"
)

var Log = logrus.New()

type ComboEdit struct {
	FileContent []byte
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
	Entries_Len int
	Entries     []Entry
}

type Entry interface {
	Name() string
	Combos() []Band
}

type UplinkEntry struct {
	combos []Band
}

func (u* UplinkEntry) Combos() []Band {
	return u.combos;
}

func (u* UplinkEntry) Name() string {
	return "UL";
}

type DownlinkEntry struct {
	combos []Band
}

func (d* DownlinkEntry) Combos() []Band {
	return d.combos;
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
	len := r.rb()
	cf.Entries_Len = int(len)
	Log.Info("This CA Combos files contains ", len, " Entries_Len")
	r.expect(0x00)

	for i:=0; i<cf.Entries_Len; i++ {
		cf.Entries = append(cf.Entries, c.parseEntry(&r))
	}

	return cf
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
		dlEntry.combos =  parseCombos(r)
		e = &dlEntry
	case 138:
		// UL
		ulEntry := UplinkEntry{}
		ulEntry.combos = parseCombos(r)
		e = &ulEntry
	default:
		Log.Warnf("Invalid type %d found!", entryType)
	}


	return e
}

func parseCombos(r *MyReader) []Band {
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