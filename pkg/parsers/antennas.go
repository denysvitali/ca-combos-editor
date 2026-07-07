package parsers

import (
	"github.com/denysvitali/ca-combos-editor/pkg/readers"
	"github.com/denysvitali/ca-combos-editor/pkg/types"
)

// ParseAntennas reads the antenna port list for a single band. The 333/334
// format stores eight raw bytes per band; non-zero bytes are interpreted as
// active antenna indices.
func ParseAntennas(r *readers.BinaryReader) []types.Antenna {
	var antennas []types.Antenna
	for range types.AntennaCount {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		if b != 0 {
			antennas = append(antennas, types.Antenna(b))
		}
	}
	return antennas
}
