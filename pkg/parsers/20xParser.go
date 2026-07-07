package parsers

import (
	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
)

// Parse20xBands reads six 4-byte band records (uint16 LE band + uint8 class +
// uint8 MIMO layer count).
func Parse20xBands(r *readers.BinaryReader) ([]types.Band, error) {
	return readBandSlots(func() (types.Band, error) {
		b, err := readBandAndClass(r)
		if err != nil {
			return b, err
		}
		mimoByte, err := r.ReadByte()
		if err != nil {
			return b, err
		}
		b.Mimo = int(mimoByte)
		return b, nil
	})
}

// Parse201 parses a downlink entry that includes MIMO information.
func Parse201(r *readers.BinaryReader) (types.Entry, error) {
	bands, err := Parse20xBands(r)
	if err != nil {
		return nil, err
	}
	return &types.DownlinkEntry{BandArr: bands}, nil
}

// Parse202 parses an uplink entry that includes MIMO information.
func Parse202(r *readers.BinaryReader) (types.Entry, error) {
	bands, err := Parse20xBands(r)
	if err != nil {
		return nil, err
	}
	return &types.UplinkEntry{BandArr: bands}, nil
}
