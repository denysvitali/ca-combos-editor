package types

import (
	"sort"
	"strings"
)

// UplinkEntry is an uplink carrier aggregation combo record.
type UplinkEntry struct {
	BandArr []Band
}

// Bands returns the bands that make up this combo.
func (u *UplinkEntry) Bands() []Band {
	return u.BandArr
}

// Name returns the entry type label.
func (u *UplinkEntry) Name() string {
	return "UL"
}

// String formats the combo as a dash-separated list of bands sorted from
// highest to lowest band number, then from highest to lowest class.
func (u *UplinkEntry) String() string {
	return joinBands(u.BandArr)
}

// SetBands replaces the bands for this combo.
func (u *UplinkEntry) SetBands(bands []Band) {
	u.BandArr = bands
}

// UlArr is a slice of uplink entries sorted by the first band number.
type UlArr []UplinkEntry

func (u UlArr) Len() int {
	return len(u)
}

// Less reports whether the i-th entry should sort before the j-th entry.
// Empty entries are placed at the end.
func (u UlArr) Less(i, j int) bool {
	if len(u[i].BandArr) == 0 {
		return false
	}
	if len(u[j].BandArr) == 0 {
		return true
	}
	return u[i].BandArr[0].Band < u[j].BandArr[0].Band
}

func (u UlArr) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
