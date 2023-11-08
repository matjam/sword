package mapgen

import "github.com/matjam/sword/internal/terrain"

////////////////////////////////////////////////////////////////////////////////
// Corridors

func (mg *MapGenerator) generateMazes() {
	// The generateMaze() method is where we generate the corridors. We do this by
	// running a maze generator at the given point.
	//
	// Maze generators work by incrementally carving passages while avoiding cutting
	// into an already open area. That's how you ensure the maze only has one solution.
	// If you let it carve into existing passages, you'd get loops.
	//
	// We deliberately do not use a recursive implementation of the maze generator,
	// because we don't want to risk overflowing the stack. Instead we use an iterative
	// implementation that uses a stack data structure to keep track of the current
	// position and the positions we've visited.

	// find all of the locations in the map that are still stone. We use this to
	// determine where we can carve new corridors. We only want to carve corridors
	// in stone tiles, not in existing corridors or rooms. We skip every second
	// tile because we only want to start a corridor at the center of each wall,
	// not at the corners.

	if mg.walking {
		mg.walk()
	} else {
		done := mg.carveMaze()
		if done {
			mg.phase = PhaseConnectors
		}
	}
}

func (mg *MapGenerator) carveMaze() (done bool) {
	// while there are still rows that have not been fully populated with rooms,
	// doors and corridors, keep carving.
	if len(mg.incompleteRows) > 0 {
		// for this row, we need to keep track of the columns that have not yet
		// been fully populated with rooms, doors and corridors.
		for x := 1; x < mg.Width-1; x += 2 {
			mg.incompleteCols = append(mg.incompleteCols, x)
		}

		// we process rows and columns in a random order and eliminate them from the
		// list once they are fully populated. This ensures that we don't end up with
		// a maze that is biased towards one direction.
		shuffleArray(mg.rng, mg.incompleteRows)
		shuffleArray(mg.rng, mg.incompleteCols)

		// we take the first row from the list of incomplete rows, which has been
		// shuffled. This ensures that we don't end up with a maze that is biased.
		scanY := mg.incompleteRows[0]

		for len(mg.incompleteCols) > 0 {
			scanX := mg.incompleteCols[0]

			if mg.terrainGrid.Get(scanX, scanY) == terrain.Stone {
				mg.x = scanX
				mg.y = scanY

				// we run a maze walker at the current position. This will carve a
				// meandering corridor until it exhausts all possible paths. We then
				// return to this method and continue scanning for incomplete rows and
				// columns.
				mg.startWalking()
				return false
			}
			// remove the column from the list of incomplete columns
			mg.incompleteCols = mg.incompleteCols[1:]

			// if we have no more incomplete columns, we're done with this row
		}

		// remove the row from the list of incomplete rows
		mg.incompleteRows = mg.incompleteRows[1:]
	} else {
		return true
	}

	return false
}

func (mg *MapGenerator) startWalking() {
	// create a new region for the maze
	mg.currentRegion = mg.nextRegion()

	// the current location was confirmed to be stone, so we set it to be a
	// corridor. We also set the regionID to the current regionID, so that we
	// can later flood fill the map to find all the disconnected regions.
	mg.terrainGrid.Set(mg.x, mg.y, terrain.Corridor)
	mg.regionGrid.Set(mg.x, mg.y, mg.currentRegion)

	// we keep track of all the locations we've visited while running the maze
	// generator. This is used by the maze hunt algorithm to find a previously
	// visited location that has an unvisited neighbour.
	mg.visitedMazeLocations = append(mg.visitedMazeLocations, [2]int{mg.x, mg.y})

	// we only start walking if we're not already walking; we don't want to start
	// a new walker if we're already walking.
	mg.walking = true
}

func (mg *MapGenerator) walk() {
	mg.walking = mg.mazeWalk()
	if !mg.walking {
		mg.walking = mg.mazeHunt()
	}

	if !mg.walking {
		// we've exhausted all possible paths from the current location, so we
		// return to the carveMaze() method, which will continue scanning for
		// incomplete rows and columns.

		// increment the regionID so that the next maze will have a different
		// regionID.
		mg.currentRegion = mg.nextRegion()
	}
}

func (mg *MapGenerator) mazeWalk() bool {
	// The mazeWalk() method is where we walk in a random direction. We do this by
	// picking a random direction, and checking if we can carve in that direction.
	// If we can, we carve a corridor in that direction, and then we start walking
	// from there. If we can't, we try again with a different random direction.
	// We keep doing this until we can't walk any further.

	directions := mg.shuffleDirections()

	for _, direction := range directions {
		face := direction
		if mg.canCarve(face) {
			mg.doCarve(face)

			// we keep track of all the locations we've visited while running the maze
			// generator. This is used by the maze hunt algorithm to find a previously
			// visited location that has an unvisited neighbour.
			mg.visitedMazeLocations = append(mg.visitedMazeLocations, [2]int{mg.x, mg.y})

			// we return true to indicate that we could walk in this direction.
			return true
		}
	}

	return false
}

func (mg *MapGenerator) mazeHunt() bool {
	// The mazeHunt() method is where we hunt for a previously visited location,
	// that has an unvisited neighbour, that is part of the same region. If we
	// find one, we set the current location to that location, return true, and
	// start walking from there. If we can't find one, we return false.

	// we shuffle the list of visited locations, so that we don't always start
	// hunting from the same location.
	shuffleArray(mg.rng, mg.visitedMazeLocations)

	for len(mg.visitedMazeLocations) > 0 {
		// try each position and see if we could walk from there
		mg.x = mg.visitedMazeLocations[0][0]
		mg.y = mg.visitedMazeLocations[0][1]

		directions := mg.shuffleDirections()
		for _, direction := range directions {
			face := direction
			if mg.canCarve(face) {
				mg.doCarve(face)
				return true
			}
		}

		// if we get here, we couldn't walk from any of the previously visited
		// locations, so we remove the current location from the list of visited
		// locations and try the next one.
		mg.visitedMazeLocations = mg.visitedMazeLocations[1:]
	}

	// if we get here, we couldn't find a previously visited location that has
	// an unvisited neighbour, so we return false.
	return false
}

func (mg *MapGenerator) shuffleDirections() []Direction {
	directions := []Direction{North, South, East, West}
	for i := range directions {
		j := mg.rng.Intn(i + 1)
		directions[i], directions[j] = directions[j], directions[i]
	}
	return directions
}

func (mg *MapGenerator) canCarve(direction Direction) bool {
	// The canCarve() method is where we check if we can carve in a given
	// direction. We do this by checking if the tile two tiles away in the given
	// direction is stone. If it is, we can carve in that direction.

	switch direction {
	case North:
		// check if the tile two tiles north is still in the terrainGrid
		if mg.y-2 < 0 {
			return false
		}
		return mg.terrainGrid.Get(mg.x, mg.y-2) == terrain.Stone
	case South:
		// check if the tile two tiles south is still in the terrainGrid
		if mg.y+2 >= mg.Height {
			return false
		}
		return mg.terrainGrid.Get(mg.x, mg.y+2) == terrain.Stone
	case East:
		// check if the tile two tiles east is still in the terrainGrid
		if mg.x+2 >= mg.Width {
			return false
		}
		return mg.terrainGrid.Get(mg.x+2, mg.y) == terrain.Stone
	case West:
		// check if the tile two tiles west is still in the terrainGrid
		if mg.x-2 < 0 {
			return false
		}
		return mg.terrainGrid.Get(mg.x-2, mg.y) == terrain.Stone
	}

	return false
}

func (mg *MapGenerator) doCarve(direction Direction) {
	// The doCarve() method is where we carve in a given direction. We do this by
	// setting the tile two tiles away in the given direction to the correct type,
	// and the tile one tile away in the given direction to the correct type.

	switch direction {
	case North:
		mg.terrainGrid.Set(mg.x, mg.y-1, terrain.Corridor)
		mg.terrainGrid.Set(mg.x, mg.y-2, terrain.Corridor)
		mg.regionGrid.Set(mg.x, mg.y-1, mg.currentRegion)
		mg.regionGrid.Set(mg.x, mg.y-2, mg.currentRegion)
		mg.y -= 2
	case South:
		mg.terrainGrid.Set(mg.x, mg.y+1, terrain.Corridor)
		mg.terrainGrid.Set(mg.x, mg.y+2, terrain.Corridor)
		mg.regionGrid.Set(mg.x, mg.y+1, mg.currentRegion)
		mg.regionGrid.Set(mg.x, mg.y+2, mg.currentRegion)
		mg.y += 2
	case East:
		mg.terrainGrid.Set(mg.x+1, mg.y, terrain.Corridor)
		mg.terrainGrid.Set(mg.x+2, mg.y, terrain.Corridor)
		mg.regionGrid.Set(mg.x+1, mg.y, mg.currentRegion)
		mg.regionGrid.Set(mg.x+2, mg.y, mg.currentRegion)
		mg.x += 2
	case West:
		mg.terrainGrid.Set(mg.x-1, mg.y, terrain.Corridor)
		mg.terrainGrid.Set(mg.x-2, mg.y, terrain.Corridor)
		mg.regionGrid.Set(mg.x-1, mg.y, mg.currentRegion)
		mg.regionGrid.Set(mg.x-2, mg.y, mg.currentRegion)
		mg.x -= 2
	}
}
