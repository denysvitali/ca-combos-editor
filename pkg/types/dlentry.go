package types

import (
	"sort"
	"strings"
)

type DownlinkEntry struct {
	BandArr []Band
}

func (d *DownlinkEntry) Name() string {
	return "DL"
}

func (d *DownlinkEntry) Bands() []Band {
	return d.BandArr
}

func (d *DownlinkEntry) String() string {
	sort.Sort(BandArr(d.BandArr))
	var bands []string

	for _, b := range d.BandArr {
		bands = append(bands, b.String())
	}

	return strings.Join(bands, "-")
}

func (d *DownlinkEntry) SetBands(bands []Band) {
	d.BandArr = bands
}
type DlArr []DownlinkEntry