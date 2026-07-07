package parsers

import (
	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
)

// Parse33xBands reads six 11-byte band records (uint16 LE band + uint8 class +
// eight antenna bytes).
func Parse33xBands(r *readers.BinaryReader) ([]types.Band, error) {
	return readBandSlots(func() (types.Band, error) {
		b, err := readBandAndClass(r)
		if err != nil {
			return b, err
		}
		b.Antennas = ParseAntennas(r)
		return b, nil
	})
}

// Parse333 parses a downlink entry that includes antenna information.
func Parse333(r *readers.BinaryReader) (types.Entry, error) {
	bands, err := Parse33xBands(r)
	if err != nil {
		return nil, err
	}
	return &types.DownlinkEntry{BandArr: bands}, nil
}

// Parse334 parses an uplink entry that includes antenna information.
func Parse334(r *readers.BinaryReader) (types.Entry, error) {
	bands, err := Parse33xBands(r)
	if err != nil {
		return nil, err
	}
	return &types.UplinkEntry{BandArr: bands}, nil
}
