package types

// EntryType identifies the kind of record stored in an NV ITEM 00028874 file.
// The low byte of the little-endian uint16 is the discriminator; the high byte
// is always zero.
type EntryType uint16

const (
	// Downlink and uplink entries without explicit MIMO information.
	EntryTypeDownlinkNoMIMO EntryType = 137
	EntryTypeUplinkNoMIMO   EntryType = 138

	// Downlink and uplink entries that also carry a MIMO layer count.
	EntryTypeDownlinkMIMO EntryType = 201
	EntryTypeUplinkMIMO   EntryType = 202

	// Downlink and uplink entries that carry an antenna port list.
	EntryTypeDownlinkAntennas EntryType = 333
	EntryTypeUplinkAntennas   EntryType = 334
)

// Binary layout constants shared by all supported entry formats.
const (
	// MaxBandsPerEntry is the number of band slots reserved for each entry.
	MaxBandsPerEntry = 6

	// HeaderSize is the size of the uncompressed payload header in bytes.
	// It consists of two zero bytes followed by a little-endian uint16 entry count.
	HeaderSize = 4

	// BandRecordSize13x is the size of a band record in 137/138 entries.
	BandRecordSize13x = 3

	// BandRecordSize20x is the size of a band record in 201/202 entries.
	BandRecordSize20x = 4

	// BandRecordSize33x is the size of a band record in 333/334 entries.
	BandRecordSize33x = 11

	// AntennaCount is the number of antenna bytes stored per band in 333/334 entries.
	AntennaCount = 8
)

// BandwidthClass maps the binary class value to the 3GPP string representation.
// A = 1, B = 2, C = 3, D = 4, E = 5. Values F-I exist in the on-air spec but
// are not commonly used in Qualcomm NV items.
type BandwidthClass int

const (
	ClassA BandwidthClass = 1
	ClassB BandwidthClass = 2
	ClassC BandwidthClass = 3
	ClassD BandwidthClass = 4
	ClassE BandwidthClass = 5
	ClassF BandwidthClass = 6
	ClassG BandwidthClass = 7
	ClassH BandwidthClass = 8
	ClassI BandwidthClass = 9
)

// MIMO layer counts commonly observed in the field.
const (
	MIMO1Layer = 1
	MIMO2Layer = 2
	MIMO4Layer = 4
)

// Valid reports whether the entry type is one of the known discriminators.
func (e EntryType) Valid() bool {
	switch e {
	case EntryTypeDownlinkNoMIMO, EntryTypeUplinkNoMIMO,
		EntryTypeDownlinkMIMO, EntryTypeUplinkMIMO,
		EntryTypeDownlinkAntennas, EntryTypeUplinkAntennas:
		return true
	}
	return false
}

// IsDownlink reports whether the entry type is a downlink record.
func (e EntryType) IsDownlink() bool {
	switch e {
	case EntryTypeDownlinkNoMIMO, EntryTypeDownlinkMIMO, EntryTypeDownlinkAntennas:
		return true
	}
	return false
}
