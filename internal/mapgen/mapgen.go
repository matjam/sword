package mapgen

import (
	"fmt"
	"image/color"
	"log/slog"
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

	Region *Region
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

type GenerationPhase int

const (
	PhaseRooms GenerationPhase = iota
	PhaseMazes
	PhaseConnectors
	PhaseConnectingRegions
	PhaseDone
)

type MapGenerator struct {
	Width  int
	Height int

	phase GenerationPhase

	maxRoomAttempts int
	curRoomAttempts int

	terrainGrid   *terrain.Terrain
	connectorGrid *grid.Grid[bool]

	roomList         []*Room
	unconnectedRooms []*Room

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

	walking    bool
	connecting bool

	rng *rand.Rand

	regions            []*Region
	unconnectedRegions []*Region
	connectedRegions   []*Region
	regionGrid         *grid.Grid[*Region]
	currentRegion      *Region
	rootRegion         *Region
}

func NewMapGenerator(width int, height int, seed int64, attempts int) *MapGenerator {
	mg := &MapGenerator{
		phase:                PhaseRooms,
		Width:                width,
		Height:               height,
		maxRoomAttempts:      attempts,
		curRoomAttempts:      0,
		terrainGrid:          terrain.NewTerrain(width, height),
		regionGrid:           grid.NewGrid[*Region](width, height),
		connectorGrid:        grid.NewGrid[bool](width, height),
		roomList:             make([]*Room, 0),
		unconnectedRooms:     make([]*Room, 0),
		incompleteRows:       make([]int, 0),
		incompleteCols:       make([]int, 0),
		visitedMazeLocations: make([][2]int, 0),
		regions:              make([]*Region, 0),
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
	//
	// This function is intended to be called in the Update() method of a game
	// loop. It will generate the map incrementally, so that you can draw the
	// map as it is being generated.

	for mg.phase != PhaseDone {
		switch mg.phase {
		case PhaseRooms:
			mg.generateRooms()
		case PhaseMazes:
			mg.generateMazes()
		case PhaseConnectors:
			mg.generateConnectors()
			return
		case PhaseConnectingRegions:
			mg.connectRegions()
			return
		default:
			return
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// Room generation

func (mg *MapGenerator) generateRooms() {
	// The generateRooms() method is where we generate the rooms. We do this by
	// picking a random room size and position, and checking if it fits. If it
	// does, we add it to the map and continue. If it doesn't, we try again with
	// a different random room size and position. We keep doing this until we
	// can't fit any more rooms into the map.

	successfullyPlacedRoom := false

	if mg.curRoomAttempts < mg.maxRoomAttempts {

		mg.currentRegion = mg.nextRegion()

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

			// We create a new room with the random size and position.
			room = Room{
				X:      roomX,
				Y:      roomY,
				Width:  roomWidth,
				Height: roomHeight,
				Region: mg.currentRegion,
			}

			// We check if the room fits in the map.
			if mg.roomFits(room) {
				mg.addRoom(room)

				successfullyPlacedRoom = true
			}

			mg.curRoomAttempts++
		}
	}

	if mg.curRoomAttempts >= mg.maxRoomAttempts {
		mg.phase = PhaseMazes
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
	for _, r := range mg.roomList {
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

	// We add the room to the list of rooms.
	mg.roomList = append(mg.roomList, &room)

	// Update the regions
	mg.regionGrid.SetRect(room.X, room.Y, room.Width, room.Height, room.Region)
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

////////////////////////////////////////////////////////////////////////////////
// Connectors

func (mg *MapGenerator) generateConnectors() {
	// The generateConnectors() method is where we generate the connectors. We do
	// this by finding all the tiles that are adjacent to a corridor, and then
	// checking if they are adjacent to a room. If they are, we add them to the
	// list of connectors. We then shuffle the list of connectors, and then we
	// iterate over the list of connectors and try to connect them to a room.

	for y := 1; y < mg.Height-1; y += 1 {
		for x := 1; x < mg.Width-1; x += 1 {
			ok, _, _ := mg.isConnectorTile(x, y)
			if ok {
				mg.connectorGrid.Set(x, y, true)
			}
		}
	}

	mg.phase = PhaseConnectingRegions
}

func (mg *MapGenerator) isConnectorTile(x, y int) (isConnector bool, region1, region2 int) {
	// Determine if the current tile connects two different regions. We only
	// conside tiles that are rooms or corridors.

	e := mg.terrainGrid.Get(x+1, y)
	w := mg.terrainGrid.Get(x-1, y)

	if (e == terrain.Room && w == terrain.Room) ||
		(e == terrain.Corridor && w == terrain.Corridor) ||
		(e == terrain.Room && w == terrain.Corridor) ||
		(e == terrain.Corridor && w == terrain.Room) {
		eRegion := mg.regionGrid.Get(x+1, y)
		wRegion := mg.regionGrid.Get(x-1, y)
		if eRegion.id != wRegion.id {
			return true, eRegion.id, wRegion.id
		}
	}

	n := mg.terrainGrid.Get(x, y-1)
	s := mg.terrainGrid.Get(x, y+1)

	if (n == terrain.Room && s == terrain.Room) ||
		(n == terrain.Corridor && s == terrain.Corridor) ||
		(n == terrain.Room && s == terrain.Corridor) ||
		(n == terrain.Corridor && s == terrain.Room) {
		nRegion := mg.regionGrid.Get(x, y-1)
		sRegion := mg.regionGrid.Get(x, y+1)
		if nRegion.id != sRegion.id {
			return true, nRegion.id, sRegion.id
		}
	}

	return false, 0, 0
}

////////////////////////////////////////////////////////////////////////////////
// Connecting regions

func (mg *MapGenerator) connectRegions() {
	if len(mg.unconnectedRegions) == 0 {
		mg.phase = PhaseDone
		return
	}

	slog.Info(fmt.Sprintf("there are %d unconnected regions", len(mg.unconnectedRegions)))
	slog.Info(fmt.Sprintf("there are %v rooms", len(mg.roomList)))

	if mg.rootRegion == nil {
		// all rooms start out as unconnected
		for _, room := range mg.roomList {
			mg.unconnectedRooms = append(mg.unconnectedRooms, room)
		}

		// shuffle the unconnected regions
		shuffleArray(mg.rng, mg.unconnectedRooms)

		// grab the last room from the list
		rootRoom := mg.unconnectedRooms[len(mg.unconnectedRooms)-1]
		mg.unconnectedRooms = mg.unconnectedRooms[:len(mg.unconnectedRooms)-1]
		mg.rootRegion = rootRoom.Region

		// remove the root region from the list of unconnected regions
		var newUnconnectedRegions []*Region
		for i, r := range mg.unconnectedRegions {
			if r.id == mg.rootRegion.id {
				newUnconnectedRegions = removeIndex(mg.unconnectedRegions, i)
				break
			}
		}
		mg.unconnectedRegions = newUnconnectedRegions

		// set the color of the root region to white
		mg.rootRegion.clr = color.RGBA{0xff, 0xff, 0xff, 0xff}

		slog.Info(fmt.Sprintf("room at %v,%v selected as root region", rootRoom.X, rootRoom.Y))

		mg.phase = PhaseDone
	}
}

////////////////////////////////////////////////////////////////////////////////
// Drawing

func (mg *MapGenerator) DrawDebug(screen *ebiten.Image) {
	for y := 0; y < mg.Height; y++ {
		for x := 0; x < mg.Width; x++ {
			t := mg.terrainGrid.Get(x, y)
			r := mg.regionGrid.Get(x, y)

			clr := color.Color(color.RGBA{0x00, 0x00, 0x00, 0xff})
			if r != nil {
				clr = r.clr
			}

			switch t {
			case terrain.Stone:
				mg.drawTile(screen, x, y, clr)
			case terrain.Room:
				mg.drawTile(screen, x, y, clr)
			case terrain.Corridor:
				mg.drawTile(screen, x, y, clr)
			case terrain.Door:
				mg.drawTile(screen, x, y, clr)
			}

			c := mg.connectorGrid.Get(x, y)
			if c {
				mg.drawDot(screen, x, y, color.RGBA{0xff, 0xff, 0xff, 0xff})
			}
		}
	}
}

func (mg *MapGenerator) drawTile(screen *ebiten.Image, x int, y int, clr color.Color) {
	vector.DrawFilledRect(screen, float32(x*16), float32(y*16), float32(16), float32(16), clr, false)
}

func (mg *MapGenerator) drawDot(screen *ebiten.Image, x int, y int, clr color.Color) {
	vector.DrawFilledRect(screen, float32(x*16+6), float32(y*16+6), float32(4), float32(4), clr, false)
}

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
		id: len(mg.regions),
		clr: color.RGBA{
			uint8(mg.rng.Intn(192) + 16),
			uint8(mg.rng.Intn(192) + 16),
			uint8(mg.rng.Intn(192) + 16),
			0xff,
		},
	}

	mg.regions = append(mg.regions, &r)
	mg.unconnectedRegions = append(mg.unconnectedRegions, &r)
	return &r
}

func removeIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
