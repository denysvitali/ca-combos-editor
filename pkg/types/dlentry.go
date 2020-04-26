package types

import (
	"sort"
	"strings"
)

type DownlinkEntry struct {
	bands []Band
}

func (d *DownlinkEntry) Name() string {
	return "DL"
}

func (d *DownlinkEntry) Bands() []Band {
	return d.bands
}

func (d *DownlinkEntry) String() string {
	sort.Sort(BandArr(d.bands))
	var bands []string

	for _, b := range d.bands {
		bands = append(bands, b.String())
	}

	return strings.Join(bands, "-")
}

func (d *DownlinkEntry) SetBands(bands []Band) {
	d.bands = bands
}
type DlArr []DownlinkEntry