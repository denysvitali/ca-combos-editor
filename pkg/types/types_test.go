package types

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBandValid(t *testing.T) {
	tests := []struct {
		name string
		band Band
		want bool
	}{
		{"valid", Band{Band: 3, Class: 1}, true},
		{"band too low", Band{Band: 0, Class: 1}, false},
		{"band too high", Band{Band: 256, Class: 1}, false},
		{"class too low", Band{Band: 3, Class: 0}, false},
		{"class too high", Band{Band: 3, Class: 10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.band.Valid())
		})
	}
}

func TestBandString(t *testing.T) {
	tests := []struct {
		band Band
		want string
	}{
		{Band{Band: 3, Class: 1, Mimo: 2}, "3A2"},
		{Band{Band: 3, Class: 1, Mimo: 1}, "3A"},
		{Band{Band: 3, Class: 6, Mimo: 4}, "3F4"},
		{Band{Band: 0, Class: 1}, ""},
		{Band{Band: 3, Class: 10}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.band.String())
		})
	}
}

func TestSortBands(t *testing.T) {
	bands := []Band{
		{Band: 3, Class: 1},
		{Band: 7, Class: 2},
		{Band: 3, Class: 2},
	}

	asc := make([]Band, len(bands))
	copy(asc, bands)
	SortBandsAsc(asc)
	assert.Equal(t, []Band{
		{Band: 3, Class: 1},
		{Band: 3, Class: 2},
		{Band: 7, Class: 2},
	}, asc)

	desc := make([]Band, len(bands))
	copy(desc, bands)
	SortBandsDesc(desc)
	assert.Equal(t, []Band{
		{Band: 7, Class: 2},
		{Band: 3, Class: 2},
		{Band: 3, Class: 1},
	}, desc)
}

func TestBandArrSort(t *testing.T) {
	bands := BandArr{
		{Band: 3, Class: 1},
		{Band: 7, Class: 2},
		{Band: 3, Class: 2},
	}
	sort.Sort(bands)
	want := BandArr{
		{Band: 7, Class: 2},
		{Band: 3, Class: 2},
		{Band: 3, Class: 1},
	}
	assert.Equal(t, want, bands)
}

func TestEntryTypeMethods(t *testing.T) {
	tests := []struct {
		et     EntryType
		valid  bool
		downlink bool
	}{
		{EntryTypeDownlinkNoMIMO, true, true},
		{EntryTypeUplinkNoMIMO, true, false},
		{EntryTypeDownlinkMIMO, true, true},
		{EntryTypeUplinkMIMO, true, false},
		{EntryTypeDownlinkAntennas, true, true},
		{EntryTypeUplinkAntennas, true, false},
		{EntryType(0), false, false},
		{EntryType(999), false, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.et), func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.et.Valid())
			assert.Equal(t, tt.downlink, tt.et.IsDownlink())
		})
	}
}
