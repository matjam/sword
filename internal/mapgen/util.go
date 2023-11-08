package mapgen

import (
	"image/color"
	"math/rand"
)

////////////////////////////////////////////////////////////////////////////////
// Utility functions

func shuffleArray[T any](rng *rand.Rand, a []T) []T {
	// woo, Fisher-Yates shuffle with generics!
	for i := len(a) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func (mg *MapGenerator) nextRegion() *Region {
	r := Region{
		id: mg.curRegionID,
		clr: color.RGBA{
			uint8(mg.rng.Intn(192) + 16),
			uint8(mg.rng.Intn(192) + 16),
			uint8(mg.rng.Intn(192) + 16),
			0xff,
		},
	}

	mg.curRegionID++
	mg.regions[r.id] = &r
	return &r
}

func removeIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
