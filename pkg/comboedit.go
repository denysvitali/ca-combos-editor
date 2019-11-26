package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var Log = logrus.New()

type ComboEdit struct {
	FileContent []byte
}

type ComboWriter struct {
	FileBody      []byte
	EntriesLength int
}

func NewComboEdit(input []byte) ComboEdit {
	return ComboEdit{FileContent: input}
}

type MyReader struct {
	reader *bytes.Reader
}

func NewMyReader(reader *bytes.Reader) MyReader {
	return MyReader{reader: reader}
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
			found&0xFF,
			b&0xFF)
	}
}

type ComboFile struct {
	EntriesLen uint16
	Entries    []Entry
}

type Entry interface {
	Name() string
	Bands() []Band
	String() string
}

type UplinkEntry struct {
	bands []Band
}

func (u *UplinkEntry) Bands() []Band {
	return u.bands
}

func (u *UplinkEntry) Name() string {
	return "UL"
}

func (u *UplinkEntry) String() string {
	sort.Sort(BandArr(u.bands))
	var bands []string

	for _, b := range u.bands {
		bands = append(bands, b.String())
	}

	return strings.Join(bands, "-")
}

type DownlinkEntry struct {
	bands []Band
}

func (d *DownlinkEntry) Bands() []Band {
	return d.bands
}

func (d *DownlinkEntry) Name() string {
	return "DL"
}

func (d *DownlinkEntry) String() string {
	sort.Sort(BandArr(d.bands))
	var bands []string

	for _, b := range d.bands {
		bands = append(bands, b.String())
	}

	return strings.Join(bands, "-")
}

type Band struct {
	Band  int
	Class int
	Mimo  int
}

func (b Band) String() string {
	// http://niviuk.free.fr/lte_ca_band.php#lte_ca_class
	// G and H are missing, TODO: inspect if QPST treats I as 7 or 9
	bandClasses := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	if b.Band < 1 || b.Band > 255 || b.Class < 1 || b.Class > 9 {
		return ""
	}
	return fmt.Sprintf("%d%s", b.Band, bandClasses[b.Class-1])
}

type BandArr []Band
type UlArr []UplinkEntry

func (u UlArr) Len() int {
	return len(u)
}

func (u UlArr) Less(i, j int) bool {
	return u[i].bands[0].Band < u[j].bands[0].Band;
}

func (u UlArr) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type DlArr []DownlinkEntry

func (b BandArr) Len() int {
	return len(b)
}

func (b BandArr) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BandArr) Less(i, j int) bool {
	return b[i].Band > b[j].Band
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
	cf.EntriesLen = binary.LittleEndian.Uint16(lenArr)
	Log.Info("This CA Bands files contains ", cf.EntriesLen, " entries")

	for i := uint16(0); i < cf.EntriesLen; i++ {
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
		fmt.Printf("Entry %s: %v\n", e.Name(), e)
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
		dlEntry.bands = parseBands(r)
		e = &dlEntry
	case 138:
		// UL
		ulEntry := UplinkEntry{}
		ulEntry.bands = parseBands(r)
		e = &ulEntry

	case 201:
		dlEntry := DownlinkEntry{}
		dlEntry.bands = parse20xBands(r)
		e = &dlEntry

	case 202:
		ulEntry := UplinkEntry{}
		ulEntry.bands = parse20xBands(r)
		e = &ulEntry
	default:
		Log.Warnf("Invalid type %d found!", entryType)
	}
	return e
}

func (c *ComboWriter) writeEntry(entry Entry) {
	switch entry.(type) {
	case *DownlinkEntry:
		c.FileBody = append(c.FileBody, byte(137))
		c.FileBody = append(c.FileBody, 0)
		sortedBands := entry.Bands()
		sort.Sort(BandArr(sortedBands))

		for i := len(sortedBands)/2-1; i >= 0; i-- {
			opp := len(sortedBands)-1-i
			sortedBands[i], sortedBands[opp] = sortedBands[opp], sortedBands[i]
		}

		c.writeBands(sortedBands)

	case *UplinkEntry:
		c.FileBody = append(c.FileBody, byte(138))
		c.FileBody = append(c.FileBody, 0)

		c.writeBands(entry.Bands())
	}
}

func parseBands(r *MyReader) []Band {
	var combos []Band

	for i := 0; i < 6; i++ {
		bwc := Band{}
		band := binary.LittleEndian.Uint16([]byte{r.rb(), r.rb()})
		class := int(r.rb())

		bwc.Band = int(band)
		bwc.Class = class

		if bwc.Band < 1 || bwc.Band > 255 || bwc.Class < 0 || bwc.Class > 9 {
			// Null, skip
			continue
		}

		combos = append(combos, bwc)
	}
	return combos
}

func parse20xBands(r *MyReader) []Band {
	var combos []Band

	for i := 0; i < 6; i++ {
		bwc := Band{}
		band := binary.LittleEndian.Uint16([]byte{r.rb(), r.rb()})
		class := int(r.rb())
		mimo := int(r.rb())

		bwc.Band = int(band)
		bwc.Class = class
		bwc.Mimo = mimo

		if bwc.Band < 1 || bwc.Band > 255 || bwc.Class < 0 || bwc.Class > 9 {
			// Null, skip
			continue
		}

		combos = append(combos, bwc)
	}
	return combos
}

func (c *ComboWriter) writeBands(bands []Band) {
	for i := 0; i < 6; i++ {
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
