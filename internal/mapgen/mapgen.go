package mapgen

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/matjam/sword/internal/grid"
	"github.com/matjam/sword/internal/terrain"
)

// roomSizes is a list of all the possible room sizes. We use this to pick a
// random room size when generating rooms. Each room size is a width and a
// height. The room sizes are all odd numbers, so that they can be centered
// inside the map. At some point we will want to support irregularly shaped
// rooms, but for now we only support rectangular rooms.
var roomSizes = [][]int{
	{3, 3}, {3, 5}, {5, 3}, {7, 5}, {5, 7}, {7, 7},
	{9, 7}, {7, 9}, {9, 9}, {11, 9}, {9, 11}, {11, 11},
}

// Room is a room in a map. They must be rectangular, and they must not overlap.
// Each room can be between 3x3 and 11x11 tiles in size. Rooms are always an odd
// number of tiles in size, so that they can be centered on the map.
// Each room's location must have an odd x and y coordinate, so that the room
// can be centered on the map.
type Room struct {
	X      int
	Y      int
	Width  int
	Height int

	RegionID int
}

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

type Region struct {
	id  int
	clr color.Color
}

type MapGenerator struct {
	Width  int
	Height int

	doneRooms      bool
	doneMazes      bool
	doneConnectors bool

	maxRoomAttempts int
	curRoomAttempts int

	terrainGrid   *terrain.Terrain
	regionGrid    *grid.Grid[Region]
	connectorGrid *grid.Grid[bool]

	roomList            []*Room
	unconnectedRoomList []*Room

	// state for maze generator
	x int
	y int

	// rows that have not yet been fully populated
	incompleteRows []int

	// when processing a row, we need to keep track of the
	// columns that have not yet been fully populated for that row
	incompleteCols []int

	// we use this to track all of the locations that we have visited
	// while running the maze generator. This is used by the maze hunt
	// algorithm to find a previously visited location that has an
	// unvisited neighbour.
	visitedMazeLocations [][2]int

	walking bool

	rng *rand.Rand

	currentRegion Region
	regions       []Region
}

func NewMapGenerator(width int, height int, seed int64, attempts int) *MapGenerator {
	mg := &MapGenerator{
		Width:                width,
		Height:               height,
		maxRoomAttempts:      attempts,
		curRoomAttempts:      0,
		terrainGrid:          terrain.NewTerrain(width, height),
		regionGrid:           grid.NewGrid[Region](width, height),
		connectorGrid:        grid.NewGrid[bool](width, height),
		roomList:             make([]*Room, 0),
		unconnectedRoomList:  make([]*Room, 0),
		incompleteRows:       make([]int, 0),
		incompleteCols:       make([]int, 0),
		visitedMazeLocations: make([][2]int, 0),
		regions:              make([]Region, 0),
	}

	for y := 1; y < mg.Height-1; y += 2 {
		mg.incompleteRows = append(mg.incompleteRows, y)
	}

	mg.rng = rand.New(rand.NewSource(seed))

	return mg
}

func (mg *MapGenerator) Update() {
	// This generate algorithm uses the "rooms and corridors" method as described
	// in this article: https://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/

	// This function is intended to be called in the Update() method of a game
	// loop. It will generate the map incrementally, so that you can draw the
	// map as it is being generated.

	mg.generateRooms()
	mg.generateMazes()
}

////////////////////////////////////////////////////////////////////////////////
// Room generation

func (mg *MapGenerator) generateRooms() {
	// The generateRooms() method is where we generate the rooms. We do this by
	// picking a random room size and position, and checking if it fits. If it
	// does, we add it to the map and continue. If it doesn't, we try again with
	// a different random room size and position. We keep doing this until we
	// can't fit any more rooms into the map.

	if mg.doneRooms {
		return
	}

	successfullyPlacedRoom := false

	if mg.curRoomAttempts < mg.maxRoomAttempts {
		for !successfullyPlacedRoom {
			var room Room

			// We generate a random room size from the list of possible room sizes.
			roomSize := roomSizes[mg.rng.Intn(len(roomSizes))]
			roomWidth := roomSize[0]
			roomHeight := roomSize[1]

			// We generate a random room position between 0 and the map width/height,
			// with an odd x and y coordinate so that rooms won't end up touching each
			// other.
			roomX := 1 + mg.rng.Intn(mg.Width/2)*2
			roomY := 1 + mg.rng.Intn(mg.Height/2)*2

			//
			mg.currentRegion = mg.nextRegion()

			// We create a new room with the random size and position.
			room = Room{
				X:      roomX,
				Y:      roomY,
				Width:  roomWidth,
				Height: roomHeight,
			}

			// We check if the room fits in the map.
			if mg.roomFits(room) {
				// We create a new region for the room.
				mg.currentRegion = mg.nextRegion()
				room.RegionID = mg.currentRegion.id
				mg.addRoom(room)

				successfullyPlacedRoom = true
			}

			mg.curRoomAttempts++
		}
	}

	if mg.curRoomAttempts >= mg.maxRoomAttempts {
		mg.doneRooms = true
	}
}

func (mg *MapGenerator) roomFits(room Room) bool {
	// The roomFits() method is where we check if a room fits in the map. We do
	// this by checking if the room overlaps with any other rooms.

	// We check if the room is inside the map.
	if room.X < 1 || room.Y < 1 || room.X+room.Width > mg.Width-1 || room.Y+room.Height > mg.Height-1 {
		return false
	}

	// We check if the room overlaps with any other rooms.
	for _, r := range mg.unconnectedRoomList {
		if room.Overlaps(r) {
			return false
		}
	}

	return true
}

func (r *Room) Overlaps(other *Room) bool {
	// The overlaps() method is where we check if a room overlaps with another
	// room. We do this by checking if the rooms overlap on the x axis and the
	// y axis.

	// We check if the rooms overlap on the x axis.
	xOverlap := r.X < other.X+other.Width && r.X+r.Width > other.X

	// We check if the rooms overlap on the y axis.
	yOverlap := r.Y < other.Y+other.Height && r.Y+r.Height > other.Y

	// We have an overlap if the rooms overlap on both axes.
	return xOverlap && yOverlap
}

func (mg *MapGenerator) addRoom(room Room) {
	// The addRoom() method is where we add a room to the map. We do this by
	// setting the tiles in the room to the correct type.
	mg.terrainGrid.SetRect(room.X, room.Y, room.Width, room.Height, terrain.Room)

	// We add the room to the list of unconnected rooms.
	mg.unconnectedRoomList = append(mg.unconnectedRoomList, &room)

	// Update the regions
	mg.regionGrid.SetRect(room.X, room.Y, room.Width, room.Height, mg.regions[room.RegionID])
}

func (mg *MapGenerator) Print() {
	// The print() method is where we print the map to the console.
	for y := 0; y < mg.Height; y++ {
		for x := 0; x < mg.Width; x++ {
			t := mg.terrainGrid.Get(x, y)
			switch t {
			case terrain.Stone:
				print("██")
			case terrain.Room:
				print("░░")
			case terrain.Corridor:
				print("  ")
			case terrain.Door:
				print("++")
			}
		}
		println()
	}
}

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

	if !mg.doneRooms || mg.doneMazes {
		return
	}

	if mg.walking {
		mg.walk()
	} else {
		mg.doneMazes = mg.carveMaze()
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
		mg.shuffleArray(mg.incompleteRows)
		mg.shuffleArray(mg.incompleteCols)

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
	// the current location was confirmed to be stone, so we set it to be a
	// corridor. We also set the regionID to the current regionID, so that we
	// can later flood fill the map to find all the disconnected regions.
	mg.terrainGrid.Set(mg.x, mg.y, terrain.Corridor)
	mg.regionGrid.Set(mg.x, mg.y, mg.currentRegion)

	// we keep track of all the locations we've visited while running the maze
	// generator. This is used by the maze hunt algorithm to find a previously
	// visited location that has an unvisited neighbour.
	mg.visitedMazeLocations = append(mg.visitedMazeLocations, [2]int{mg.x, mg.y})

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
	mg.shufflePositionArray(mg.visitedMazeLocations)

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

func (mg *MapGenerator) DrawDebug(screen *ebiten.Image) {
	for y := 0; y < mg.Height; y++ {
		for x := 0; x < mg.Width; x++ {
			t := mg.terrainGrid.Get(x, y)
			switch t {
			case terrain.Stone:
				mg.drawTile(screen, x, y, color.RGBA{0x33, 0x33, 0x33, 0xff})
			case terrain.Room:
				mg.drawTile(screen, x, y, color.RGBA{0x00, 0x99, 0x00, 0xff})
			case terrain.Corridor:
				mg.drawTile(screen, x, y, color.RGBA{0x99, 0x99, 0x99, 0xff})
			case terrain.Door:
				mg.drawTile(screen, x, y, color.RGBA{0x99, 0x00, 0x00, 0xff})
			}
		}
	}
}

func (mg *MapGenerator) drawTile(screen *ebiten.Image, x int, y int, clr color.Color) {
	vector.DrawFilledRect(screen, float32(x*16), float32(y*16), float32(16), float32(16), clr, false)
}

////////////////////////////////////////////////////////////////////////////////
// Utility functions

func (mg *MapGenerator) shuffleArray(a []int) {
	for i := len(a) - 1; i > 0; i-- {
		j := mg.rng.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return
}

func (mg *MapGenerator) shufflePositionArray(a [][2]int) {
	for i := len(a) - 1; i > 0; i-- {
		j := mg.rng.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return
}

func (mg *MapGenerator) nextRegion() Region {
	r := Region{
		id: len(mg.regions),
		clr: color.RGBA{
			uint8(mg.rng.Intn(255)),
			uint8(mg.rng.Intn(255)),
			uint8(mg.rng.Intn(255)),
			0xff,
		},
	}

	mg.regions = append(mg.regions, r)
	return r
}
