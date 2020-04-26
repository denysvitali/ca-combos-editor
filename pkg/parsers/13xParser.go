package parsers

import (
	"encoding/binary"
	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
)

func Parse13xBands(r *readers.BinaryReader) []types.Band {
	var combos []types.Band

	for i := 0; i < 6; i++ {
		bwc := types.Band{}
		band := binary.LittleEndian.Uint16([]byte{r.Rb(), r.Rb()})
		class := int(r.Rb())
		bwc.Band = int(band)
		bwc.Class = class

		if bwc.Band < 1 || bwc.Band > 255 || bwc.Class < 0 || bwc.Class > 9 {
			// Null, skip
			continue
		}
		combos = append(combos, bwc)
	}
	return combos
}

func Parse137(r *readers.BinaryReader) types.Entry {
	entry := &types.DownlinkEntry{}
	entry.SetBands(Parse13xBands(r))
	return entry
}

func Parse138(r *readers.BinaryReader) types.Entry {
	entry := &types.UplinkEntry{}
	entry.SetBands(Parse13xBands(r))
	return entry
}