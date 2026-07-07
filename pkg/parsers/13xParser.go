package parsers

import (
	"encoding/binary"
	"fmt"

	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
)

// bandSlotReader parses one band slot and returns the populated Band.
type bandSlotReader func() (types.Band, error)

// readBandSlots reads exactly MaxBandsPerEntry band slots using slot for the
// per-slot layout. Slots with an invalid band or class are treated as empty and
// skipped.
func readBandSlots(slot bandSlotReader) ([]types.Band, error) {
	var combos []types.Band
	for i := range types.MaxBandsPerEntry {
		b, err := slot()
		if err != nil {
			return nil, fmt.Errorf("band slot %d: %w", i, err)
		}
		if b.Band >= 1 && b.Band <= 255 && b.Class >= 1 && b.Class <= 9 {
			combos = append(combos, b)
		}
	}
	return combos, nil
}

// Parse13xBands reads six 3-byte band records (uint16 LE band + uint8 class).
func Parse13xBands(r *readers.BinaryReader) ([]types.Band, error) {
	return readBandSlots(func() (types.Band, error) {
		b := types.Band{}
		bandBytes, err := r.ReadBytes(2)
		if err != nil {
			return b, err
		}
		b.Band = int(binary.LittleEndian.Uint16(bandBytes))
		classByte, err := r.Rb()
		if err != nil {
			return b, err
		}
		b.Class = int(classByte)
		return b, nil
	})
}

// Parse137 parses a downlink entry without MIMO information.
func Parse137(r *readers.BinaryReader) (types.Entry, error) {
	bands, err := Parse13xBands(r)
	if err != nil {
		return nil, err
	}
	return &types.DownlinkEntry{BandArr: bands}, nil
}

// Parse138 parses an uplink entry without MIMO information.
func Parse138(r *readers.BinaryReader) (types.Entry, error) {
	bands, err := Parse13xBands(r)
	if err != nil {
		return nil, err
	}
	return &types.UplinkEntry{BandArr: bands}, nil
}
