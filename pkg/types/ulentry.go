package types

import (
	"sort"
	"strings"
)

type UplinkEntry struct {
	BandArr []Band
}

func (u *UplinkEntry) Bands() []Band {
	return u.BandArr
}

func (u *UplinkEntry) Name() string {
	return "UL"
}

func (u *UplinkEntry) String() string {
	sort.Sort(BandArr(u.BandArr))
	var bands []string

	for _, b := range u.BandArr {
		bands = append(bands, b.String())
	}

	return strings.Join(bands, "-")
}

func (u *UplinkEntry) SetBands(bands []Band){
	u.BandArr = bands
}

type UlArr []UplinkEntry

func (u UlArr) Len() int {
	return len(u)
}

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
