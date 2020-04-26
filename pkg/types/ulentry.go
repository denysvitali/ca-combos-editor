package types

import (
	"sort"
	"strings"
)

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

func (u *UplinkEntry) SetBands(bands []Band){
	u.bands = bands
}

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
