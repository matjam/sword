package terrain

import "github.com/matjam/sword/internal/grid"

// package terrain defines a terrain system for the game that we can use
// to generate the tilemap for the game, based on the rules defined in the
// terrain system.

type Type uint8

const (
	Stone Type = iota
	Room
	Corridor
	Door
)

type Terrain struct {
	*grid.Grid[Type]

	Width  int
	Height int
}

// NewTerrain creates a new terrain grid with the given width and height. The
// grid is initially filled with Stone.
func NewTerrain(width, height int) *Terrain {
	return &Terrain{
		Width:  width,
		Height: height,
		Grid:   grid.NewGrid[Type](width, height),
	}
}
