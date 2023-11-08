package mapgen

import "github.com/matjam/sword/internal/terrain"

////////////////////////////////////////////////////////////////////////////
// Remove Dead Ends

func (mg *MapGenerator) removeDeadEnds() {
	// The removeDeadEnds() method is where we remove dead ends. We do this by
	// iterating over the map, and for each tile we check if it is a dead end. If
	// it is, we remove it.

	mg.deadEndsPreviouslyRemoved = mg.deadEndsRemoved

	mg.findDeadEnds()
	for _, deadEnd := range mg.deadEnds {
		x, y := deadEnd[0], deadEnd[1]
		mg.terrainGrid.Set(x, y, terrain.Stone)
		mg.regionGrid.Set(x, y, nil)
		mg.deadEndsRemoved++
	}
	if mg.deadEndsPreviouslyRemoved == mg.deadEndsRemoved {
		mg.Phase = PhaseDone
	}
}

func (mg *MapGenerator) isDeadEnd(x, y int) bool {
	// The isDeadEnd() method is where we check if a tile is a dead end. We do
	// this by checking if the tile is a corridor, and if it has only one
	// neighbouring corridor tile.

	t := mg.terrainGrid.Get(x, y)
	if t != terrain.Corridor && t != terrain.Door {
		return false
	}

	neighbours := mg.getNeighbours(x, y)

	// count the number of corridor neighbours
	corridorNeighbours := 0
	for _, n := range neighbours {
		if n != terrain.Stone {
			corridorNeighbours++
		}
	}

	return corridorNeighbours == 1
}

func (mg *MapGenerator) getNeighbours(x, y int) []terrain.Type {
	// The getNeighbours() method is where we get the neighbours of a tile. We do
	// this by getting the tile to the north, south, east and west of the given
	// tile.

	n := mg.terrainGrid.Get(x, y-1)
	s := mg.terrainGrid.Get(x, y+1)
	e := mg.terrainGrid.Get(x+1, y)
	w := mg.terrainGrid.Get(x-1, y)

	return []terrain.Type{n, s, e, w}
}

func (mg *MapGenerator) findDeadEnds() {
	// The findDeadEnds() method is where we find all the dead ends in the map. We
	// do this by iterating over the map, and for each tile we check if it is a
	// dead end. If it is, we add it to the list of dead ends.

	mg.deadEnds = make([][2]int, 0)

	for y := 0; y < mg.Height; y++ {
		for x := 0; x < mg.Width; x++ {
			if mg.isDeadEnd(x, y) {
				mg.deadEnds = append(mg.deadEnds, [2]int{x, y})
			}
		}
	}
}
