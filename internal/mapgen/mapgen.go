package mapgen

import (
	"image/color"
	"log/slog"
	"math/rand"
	"time"

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

type RegionID int
type Region struct {
	id  RegionID
	clr color.Color
}

type Connector struct {
	x, y             int
	region1, region2 *Region
}

type GenerationPhase int

const (
	PhaseRooms GenerationPhase = iota
	PhaseMazes
	PhaseConnectors
	PhaseConnectingRegions
	PhaseRemoveDeadEnds
	PhaseDone
)

type MapGenerator struct {
	Width  int
	Height int

	phase GenerationPhase

	maxRoomAttempts int
	curRoomAttempts int

	terrainGrid   *terrain.Terrain
	connectorGrid *grid.Grid[*Connector]
	regionGrid    *grid.Grid[*Region]

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

	curRegionID   RegionID
	regions       map[RegionID]*Region
	currentRegion *Region
	rootRegion    *Region

	connectors     []*Connector
	rootConnectors []*Connector

	deadEnds                  [][2]int
	deadEndsRemoved           int
	deadEndsPreviouslyRemoved int
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
		connectorGrid:        grid.NewGrid[*Connector](width, height),
		roomList:             make([]*Room, 0),
		unconnectedRooms:     make([]*Room, 0),
		incompleteRows:       make([]int, 0),
		incompleteCols:       make([]int, 0),
		visitedMazeLocations: make([][2]int, 0),
		regions:              make(map[RegionID]*Region),
		connectors:           make([]*Connector, 0),
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

	startTime := time.Now()
	for mg.phase != PhaseDone {
		switch mg.phase {
		case PhaseRooms:
			mg.generateRooms()
		case PhaseMazes:
			mg.generateMazes()
		case PhaseConnectors:
			mg.generateConnectors()
		case PhaseConnectingRegions:
			mg.connectRegions()
		case PhaseRemoveDeadEnds:
			mg.removeDeadEnds()
		default:
			return
		}
	}
	endTime := time.Now()

	slog.Debug("Map generation finished", "time", endTime.Sub(startTime))
}

////////////////////////////////////////////////////////////////////////////////
// Remove dead ends
