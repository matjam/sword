package mapgen

import "github.com/matjam/sword/internal/terrain"

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
