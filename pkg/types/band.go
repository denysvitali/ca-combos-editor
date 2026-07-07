package types

import (
	"cmp"
	"fmt"
	"slices"
)

const (
	minBand  = 1
	maxBand  = 255
	maxClass = 9
)

// bandClasses maps the binary class value (1..9) to the 3GPP string.
// A=1, B=2, C=3, D=4, E=5, F=6, G=7, H=8, I=9.
var bandClasses = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

// Band represents one component of a carrier aggregation combo.
type Band struct {
	Band     int
	Class    int
	Mimo     int
	Antennas []Antenna
}

// Valid reports whether the band number and class are in the populated range.
func (b Band) Valid() bool {
	return b.Band >= minBand && b.Band <= maxBand &&
		b.Class >= 1 && b.Class <= maxClass
}

// ClassString returns the 3GPP bandwidth class letter for the band.
func (b Band) ClassString() string {
	if b.Class < 1 || b.Class > maxClass {
		return ""
	}
	return bandClasses[b.Class-1]
}

func (b Band) String() string {
	if !b.Valid() {
		return ""
	}
	if b.Mimo > 1 {
		return fmt.Sprintf("%d%s%d", b.Band, b.ClassString(), b.Mimo)
	}
	return fmt.Sprintf("%d%s", b.Band, b.ClassString())
}

// BandArr is a slice of bands that sorts in descending order by band number,
// then by class. It implements sort.Interface.
type BandArr []Band

// Len returns the number of bands in the slice.
func (b BandArr) Len() int { return len(b) }

// Swap exchanges the elements at indices i and j.
func (b BandArr) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// Less reports whether the band at index i should sort before the band at index j.
func (b BandArr) Less(i, j int) bool {
	if b[i].Band == b[j].Band {
		return b[i].Class > b[j].Class
	}
	return b[i].Band > b[j].Band
}

// SortBandsDesc sorts bands by descending band number, then descending class.
func SortBandsDesc(bands []Band) {
	slices.SortFunc(bands, func(a, b Band) int {
		if n := cmp.Compare(b.Band, a.Band); n != 0 {
			return n
		}
		return cmp.Compare(b.Class, a.Class)
	})
}

// SortBandsAsc sorts bands by ascending band number, then ascending class.
func SortBandsAsc(bands []Band) {
	slices.SortFunc(bands, func(a, b Band) int {
		if n := cmp.Compare(a.Band, b.Band); n != 0 {
			return n
		}
		return cmp.Compare(a.Class, b.Class)
	})
}
