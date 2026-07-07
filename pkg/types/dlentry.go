package types

import (
	"cmp"
	"slices"
	"strings"
)

// DownlinkEntry is a downlink carrier aggregation combo record.
type DownlinkEntry struct {
	BandArr []Band
}

// Name returns the entry type label.
func (d *DownlinkEntry) Name() string {
	return "DL"
}

// Bands returns the bands that make up this combo.
func (d *DownlinkEntry) Bands() []Band {
	return d.BandArr
}

// String formats the combo as a dash-separated list of bands sorted from
// highest to lowest band number, then from highest to lowest class.
func (d *DownlinkEntry) String() string {
	slices.SortFunc(d.BandArr, func(a, b Band) int {
		if a.Band != b.Band {
			return cmp.Compare(b.Band, a.Band)
		}
		return cmp.Compare(b.Class, a.Class)
	})
	var bands []string

	for _, b := range d.BandArr {
		bands = append(bands, b.String())
	}

	return strings.Join(bands, "-")
}

// SetBands replaces the bands for this combo.
func (d *DownlinkEntry) SetBands(bands []Band) {
	d.BandArr = bands
}

// DlArr is a slice of downlink entries.
type DlArr []DownlinkEntry
