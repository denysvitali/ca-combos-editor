package types

import "fmt"

// Band represents one component of a carrier aggregation combo.
type Band struct {
	Band     int
	Class    int
	Mimo     int
	Antennas []Antenna
}

// bandClasses maps the binary class value (1..9) to the 3GPP string.
// A=1, B=2, C=3, D=4, E=5, F=6, G=7, H=8, I=9.
var bandClasses = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

func (b Band) String() string {
	if b.Band < 1 || b.Band > 255 || b.Class < 1 || b.Class > len(bandClasses) {
		return ""
	}
	mimoString := ""
	if b.Mimo > 1 {
		mimoString = fmt.Sprintf("%d", b.Mimo)
	}
	return fmt.Sprintf("%d%s%s", b.Band, bandClasses[b.Class-1], mimoString)
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
