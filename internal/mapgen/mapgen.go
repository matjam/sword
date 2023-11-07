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

	Region int
}

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

type MapGenerator struct {
	Width  int
	Height int

	terrainGrid   *terrain.Terrain
	regionGrid    *grid.Grid[int]
	connectorGrid *grid.Grid[bool]

	roomList            []*Room
	unconnectedRoomList []*Room

	// state for maze generator
	x int
	y int

	rng *rand.Rand
}

func NewMapGenerator(width int, height int) *MapGenerator {
	return &MapGenerator{
		Width:               width,
		Height:              height,
		terrainGrid:         terrain.NewTerrain(width, height),
		regionGrid:          grid.NewGrid[int](width, height),
		roomList:            make([]*Room, 0),
		unconnectedRoomList: make([]*Room, 0),
	}
}

func (mg *MapGenerator) Generate(seed int64) {
	// This generate algorithm uses the "rooms and corridors" method as described
	// in this article: https://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/
	//
	// The basic idea is that we start with a blank map, and then we add rooms to
	// the map. We do this by picking a random room size and position, and
	// checking if it fits. If it does, we add it to the map and continue. If it
	// doesn't, we try again with a different random room size and position. We
	// keep doing this until we can't fit any more rooms into the map.
	//
	// Once we have all the rooms in the map, Then we iterate over every tile in the
	// dungeon. When we find a solid tile where an open area could be, we start
	// running a maze generator at that point.
	//
	// Maze generators work by incrementally carving passages while avoiding cutting
	// into an already open area. That's how you ensure the maze only has one solution.
	// If you let it carve into existing passages, you'd get loops.
	//
	// This is conveniently exactly what you need to let the maze grow and fill the
	// odd shaped areas that surround the rooms. In other words, a maze generator is
	// a randomized flood fill algorithm. Run this on every solid region between the
	// rooms and we're left with the entire dungeon packed full of disconnected rooms
	// and mazes.
	//
	// We then iterate over every tile in the dungeon again. When we find a wall tile
	// that

	mg.rng = rand.New(rand.NewSource(seed))

	mg.generateRooms(250)
	mg.generateMazes()
}

////////////////////////////////////////////////////////////////////////////////
// Room generation

func (mg *MapGenerator) generateRooms(attempts int) {
	// The generateRooms() method is where we generate the rooms. We do this by
	// picking a random room size and position, and checking if it fits. If it
	// does, we add it to the map and continue. If it doesn't, we try again with
	// a different random room size and position. We keep doing this until we
	// can't fit any more rooms into the map.

	for i := 0; i < attempts; i++ {
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

		// We create a new room with the random size and position.
		room = Room{
			X:      roomX,
			Y:      roomY,
			Width:  roomWidth,
			Height: roomHeight,
		}

		// We check if the room fits in the map.
		if mg.roomFits(room) {
			mg.addRoom(room)
		}
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
	mg.regionGrid.SetRect(room.X, room.Y, room.Width, room.Height, len(mg.roomList))
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

	mg.x = 1
	mg.y = 1

	mg.carveMaze()
}

func (mg *MapGenerator) carveMaze() {
	var x, y int
	roomExists := true
	for roomExists {
		x = 1 + mg.rng.Intn(mg.Width/2)*2
		y = 1 + mg.rng.Intn(mg.Height/2)*2
		if mg.terrainGrid.Get(x, y) != terrain.Room {
			roomExists = false
		}
	}

	// scan for walls
	for mg.hasUnconnectedStone() {
		mg.startWalking()
	}
}

func (mg *MapGenerator) hasUnconnectedStone() bool {
	// The hasUnconnectedStone() method is where we check if there are any
	// unconnected stone tiles left in the map. We use this to determine if we
	// should keep generating corridors.
	foundStone := false

	// We iterate over every tile in the map.
	for y := 1; y < mg.Height-1; y += 2 {
		for x := 1; x < mg.Width-1; x += 2 {
			// We check if the tile is stone.
			if mg.terrainGrid.Get(x, y) == terrain.Stone {
				// We check if the tile has any stone neighbours.
				foundStone = true
				mg.x = x
				mg.y = y
				return foundStone
			}
		}
	}

	return foundStone
}

func (mg *MapGenerator) startWalking() {
	mg.terrainGrid.Set(mg.x, mg.y, terrain.Corridor)

	for ok := true; ok; {
		ok = mg.mazeWalk()
		if !ok {
			ok = mg.mazeHunt()
		}
	}
}

func (mg *MapGenerator) mazeWalk() bool {
	directions := mg.shuffleDirections()

	for _, direction := range directions {
		face := direction
		if mg.canCarve(face) {
			mg.doCarve(face)
			return true
		}
	}

	return false
}

func (mg *MapGenerator) shuffleArray(a []int) {
	for i := len(a) - 1; i > 0; i-- {
		j := mg.rng.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return
}

// hunt for a previously visited location, that has an unvisited neighbour.
func (mg *MapGenerator) mazeHunt() bool {
	incompleteRows := make([]int, 0)

	for y := 1; y < mg.Height-1; y += 2 {
		incompleteRows = append(incompleteRows, y)
		mg.shuffleArray(incompleteRows)
	}

	for len(incompleteRows) > 0 {
		scanY := incompleteRows[0]
		scanX := 0

		for scanX = 1; scanX < mg.Width-1; scanX += 2 {
			if mg.terrainGrid.Get(scanX, scanY) == terrain.Corridor {
				if mg.canCarve(North) {
					mg.x = scanX
					mg.y = scanY
					mg.doCarve(North)
					return true
				}

				if mg.canCarve(South) {
					mg.x = scanX
					mg.y = scanY
					mg.doCarve(South)
					return true
				}

				if mg.canCarve(East) {
					mg.x = scanX
					mg.y = scanY
					mg.doCarve(East)
					return true
				}

				if mg.canCarve(West) {
					mg.x = scanX
					mg.y = scanY
					mg.doCarve(West)
					return true
				}
			}
		}
		// fuond a hallway but nothing I could use
		// remove the row from the list of incomplete rows
		if scanX >= mg.Width-1 {
			incompleteRows = incompleteRows[1:]
		}
	}

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
		mg.y -= 2
	case South:
		mg.terrainGrid.Set(mg.x, mg.y+1, terrain.Corridor)
		mg.terrainGrid.Set(mg.x, mg.y+2, terrain.Corridor)
		mg.y += 2
	case East:
		mg.terrainGrid.Set(mg.x+1, mg.y, terrain.Corridor)
		mg.terrainGrid.Set(mg.x+2, mg.y, terrain.Corridor)
		mg.x += 2
	case West:
		mg.terrainGrid.Set(mg.x-1, mg.y, terrain.Corridor)
		mg.terrainGrid.Set(mg.x-2, mg.y, terrain.Corridor)
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
