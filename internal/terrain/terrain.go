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

// Directions is a bitfield that represents the neighbours of a tile. It is
// used by the terrain system to determine what type of tile to place next
// to another tile.
type Directions uint8

// The bitfield values for each direction.
const (
	North     Directions = 1
	NorthEast            = 2
	East                 = 4
	SouthEast            = 8
	South                = 16
	SouthWest            = 32
	West                 = 64
	NorthWest            = 128
)

// GetNeighbors returns the neighbours of the given tile
// as a bitfield. If the tile is a wall, the bit for that
// direction will be set.
//
// The bitfield is as follows:
//
//  8 1 2
//  7   3
//  6 5 4
//
// So, for example, if the tile to the north is stone, the
// bitfield will be 0000 0001, or 1. A wall to the north and
// east would be 0000 0100, or 4, and so on.
//
// In this context, only stone tiles are considered solid.
// Doors and floors are not.
func (m *Terrain) GetNeighbors(x, y int) Directions {
	var neighbors Directions

	if m.Get(x, y) == Stone {
		neighbors |= North
	}

	if m.Get(x, y) == Stone {
		neighbors |= NorthEast
	}

	if m.Get(x, y) == Stone {
		neighbors |= East
	}

	if m.Get(x, y) == Stone {
		neighbors |= SouthEast
	}

	if m.Get(x, y) == Stone {
		neighbors |= South
	}

	if m.Get(x, y) == Stone {
		neighbors |= SouthWest
	}

	if m.Get(x, y) == Stone {
		neighbors |= West
	}

	if m.Get(x, y) == Stone {
		neighbors |= NorthWest
	}

	return neighbors
}
