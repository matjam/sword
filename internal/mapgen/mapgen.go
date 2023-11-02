package mapgen

import (
	"math/rand"

	"github.com/matjam/sword/internal/tilemap"
)

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

type MapGenerator struct {
	Width  int
	Height int

	Tilemap *tilemap.Tilemap

	roomList            []*Room
	unconnectedRoomList []*Room
}

func NewMapGenerator(width int, height int) *MapGenerator {
	return &MapGenerator{
		Width:               width,
		Height:              height,
		Tilemap:             tilemap.NewTilemap(width, height),
		roomList:            make([]*Room, 0),
		unconnectedRoomList: make([]*Room, 0),
	}
}

func (mg *MapGenerator) Generate() {
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

	mg.generateRooms()
	mg.generateCorridors()
}

func (mg *MapGenerator) generateRooms() {
	// The generateRooms() method is where we generate the rooms. We do this by
	// picking a random room size and position, and checking if it fits. If it
	// does, we add it to the map and continue. If it doesn't, we try again with
	// a different random room size and position. We keep doing this until we
	// can't fit any more rooms into the map.

	// We start with a blank tilemap.
	mg.Tilemap = tilemap.NewTilemap(mg.Width, mg.Height)

	// pick a total number of rooms based on the map size
	roomCount := 10 + rand.Intn(mg.Width*mg.Height/100)

	for i := 0; i < roomCount; i++ {
		roomFits := false
		var room Room

		for !roomFits {
			// We generate a random room size between 3x3 and 11x11.
			roomWidth := 3 + rand.Intn(9)
			roomHeight := 3 + rand.Intn(9)

			// We generate a random room position between 0 and the map width/height.
			roomX := rand.Intn(mg.Width - roomWidth)
			roomY := rand.Intn(mg.Height - roomHeight)

			// We create a new room with the random size and position.
			room = Room{
				X:      roomX,
				Y:      roomY,
				Width:  roomWidth,
				Height: roomHeight,
			}

			// We check if the room fits in the map.
			roomFits = mg.roomFits(room)
		}
		mg.addRoom(room)
	}
}

func (mg *MapGenerator) roomFits(room Room) bool {
	// The roomFits() method is where we check if a room fits in the map. We do
	// this by checking if the room overlaps with any other rooms in the map.

	// We loop through all the rooms in the map.
	for _, r := range mg.roomList {
		// We check if the room overlaps with the current room.
		if room.X < r.X+r.Width && room.X+room.Width > r.X &&
			room.Y < r.Y+r.Height && room.Y+room.Height > r.Y {
			// If it does, we return false.
			return false
		}
	}

	// If it doesn't, we return true.
	return true
}

func (mg *MapGenerator) addRoom(room Room) {
	// The addRoom() method is where we add a room to the map. We do this by
	// setting the tiles in the room to the correct type.

	// We loop through all the tiles in the room.
	for y := room.Y; y < room.Y+room.Height; y++ {
		for x := room.X; x < room.X+room.Width; x++ {
			// We set the tile at the current position to the correct type.
			mg.Tilemap.SetTile(x, y, &tilemap.Tile{
				Type: tilemap.TileTypeFloor,
			})
		}
	}

	// We add the room to the list of unconnected rooms.
	mg.unconnectedRoomList = append(mg.unconnectedRoomList, &room)
}

func (mg *MapGenerator) generateCorridors() {
	// The generateCorridors() method is where we generate the corridors. We do
	// this by iterating over every tile in the dungeon. When we find a wall
	// tile where an open area could be, we start running a maze generator at
	// that point.

	// We loop through all the tiles in the dungeon.
	for y := 1; y < mg.Height; y = y + 2 {
		for x := 1; x < mg.Width; x = x + 2 {
			// We check if the tile is solid and if it's surrounded by wall tiles.
			if mg.Tilemap.GetTile(x, y).Type == tilemap.TileTypeWall &&
				mg.isSurroundedByTileType(x, y, tilemap.TileTypeWall) {
				// If it is, we run a maze generator at that point.
				mg.generateMaze(x, y)
			}
		}
	}
}

func (mg *MapGenerator) generateMaze(x int, y int) {
	// The generateMaze() method is where we generate a maze. We do this by
	// running a maze generator at a given point.

	// pick a random direction to start with
	direction := rand.Intn(4)

	// loop until we hit a dead end
	for {
		// dig out the next two tiles in the current direction
		mg.Tilemap.SetTile(x, y, &tilemap.Tile{
			Type: tilemap.TileTypeFloor,
		})
		mg.Tilemap.SetTile(x+((direction+1)%2)*2-1, y+(direction%2)*2-1, &tilemap.Tile{
			Type: tilemap.TileTypeFloor,
		})
		mg.Tilemap.SetTile(x+((direction+1)%2)*2, y+(direction%2)*2, &tilemap.Tile{
			Type: tilemap.TileTypeFloor,
		})

		// move forward in the current direction
		x += (direction+1)%2*2 - 1
		y += direction%2*2 - 1

		// check if we hit a dead end
		if !mg.isSurroundedByTileType(x, y, tilemap.TileTypeWall) {
			// if not, pick a new direction
			direction = rand.Intn(4)
		} else {
			// if so, we're done
			break
		}
	}
}

func (mg *MapGenerator) isSurroundedByTileType(x int, y int, tt tilemap.TileType) bool {
	// The isSurroundedByTileType() method is where we check if a tile is
	// surrounded by a specific tile type. We do this by checking if the tiles
	// around the current tile are the specified tile type.

	// We loop through all the tiles around the current tile.
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// We check if the tile is the specified tile type.
			if mg.Tilemap.GetTile(x+i, y+j).Type != tt {
				// If it isn't, we return false.
				return false
			}
		}
	}

	// If it is, we return true.
	return true
}
