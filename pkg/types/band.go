package types

import "fmt"

type Band struct {
	Band     int
	Class    int
	Mimo     int
	Antennas []Antenna
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
func (b BandArr) Len() int {
	return len(b)
}

func (b BandArr) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BandArr) Less(i, j int) bool {
	return b[i].Band > b[j].Band
}