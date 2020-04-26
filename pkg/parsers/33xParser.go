package parsers

import (
	"encoding/binary"
	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
)

func Parse33xBand(r *readers.BinaryReader) []types.Band {
	var combos []types.Band

	for i := 0; i < 6; i++ {
		bwc := types.Band{}
		band := binary.LittleEndian.Uint16(r.ReadBytes(2))
		class := int(r.Rb())
		antennas := ParseAntennas(r)

		bwc.Band = int(band)
		bwc.Class = class
		bwc.Antennas = antennas

		if bwc.Band < 1 || bwc.Band > 255 || bwc.Class < 0 || bwc.Class > 9 {
			// Null, skip
			continue
		}

		combos = append(combos, bwc)
	}
	return combos
}

func Parse333(r *readers.BinaryReader) types.Entry {
	entry := &types.DownlinkEntry{}
	entry.SetBands(Parse33xBand(r))
	return entry
}

func Parse334(r *readers.BinaryReader) types.Entry {
	entry := &types.UplinkEntry{}
	entry.SetBands(Parse33xBand(r))
	return entry
}
