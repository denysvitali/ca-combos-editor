package parsers

import (
	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
)

func ParseAntennas(r *readers.BinaryReader) []types.Antenna {
	var antennas []types.Antenna
	for i := 0; i<8; i++ {
		antennaEntry := types.Antenna(r.Rb())
		if antennaEntry != 0 {
			antennas = append(antennas, antennaEntry)
		}
	}
	return antennas
}