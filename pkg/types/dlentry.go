package types

import "strings"

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
	return joinBands(d.BandArr)
}

// SetBands replaces the bands for this combo.
func (d *DownlinkEntry) SetBands(bands []Band) {
	d.BandArr = bands
}

// DlArr is a slice of downlink entries.
type DlArr []DownlinkEntry

// joinBands formats bands as a dash-separated string in descending order.
func joinBands(bands []Band) string {
	if len(bands) == 0 {
		return ""
	}
	sorted := make([]Band, len(bands))
	copy(sorted, bands)
	SortBandsDesc(sorted)

	parts := make([]string, len(sorted))
	for i, b := range sorted {
		parts[i] = b.String()
	}
	return strings.Join(parts, "-")
}
