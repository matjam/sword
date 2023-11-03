// Package tilemap implements a tilemap. It holds a grid of tiles. Rendering
// is handled by separate renderers implemented for specific target displays.
// For example, there is a renderer for the terminal, and a renderer for a
// graphical display.
package tilemap

//go:generate go-enum --marshal

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer interface {
	// Draw is called every frame to draw the grid to the screen.
	Draw(dst *ebiten.Image, x int, y int, viewport Rectangle)
}

type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

// ENUM(wall, closed_door, open_door, floor, stairs_up, stairs_down)
type TileType uint8

// Tile is a single tile in a grid. The Tile struct holds information about
// whether the tile has been seen by the player, and what region it belongs to
// which is used during map generation.
type Tile struct {
	Type       TileType
	Region     int
	Seen       bool
	Visible    bool
	LightLevel uint8
}

// Grid is a map of tiles. It holds information about the size of the map,
// and a slice of tiles. Grids do not handle any the rendering of the map,
// they only hold the data.
type Grid struct {
	Width  int
	Height int
	Tiles  []Tile
}

// NewGrid creates a new Grid with the given width and height.
func NewGrid(width int, height int) *Grid {
	tm := &Grid{
		Width:  width,
		Height: height,
		Tiles:  make([]Tile, width*height),
	}

	for i := 0; i < width*height; i++ {
		tm.Tiles[i].Type = TileTypeWall
	}
	return tm
}

// GetTile returns the tile at the given position. If the position is outside
// the bounds of the map, it returns nil.
func (tm *Grid) GetTile(x int, y int) *Tile {
	if x < 0 || x >= tm.Width || y < 0 || y >= tm.Height {
		return nil
	}
	return &tm.Tiles[y*tm.Width+x]
}

// SetTile sets the tile at the given position to the given tile. If the
// position is outside the bounds of the map, it does nothing.
func (tm *Grid) SetTile(x int, y int, tile *Tile) {
	if x < 0 || x >= tm.Width || y < 0 || y >= tm.Height {
		return
	}
	tm.Tiles[y*tm.Width+x] = *tile
}

// IsVisible returns true if the tile at the given position is visible to the
// second tile at the given position. If either of the positions are outside
// the bounds of the map, it returns false. This is calculated dynamically by
// performing a line of sight check between the two tiles.
func (tm *Grid) IsVisible(x1 int, y1 int, x2 int, y2 int) bool {
	// If either of the positions are outside the bounds of the map, we return
	// false.
	if x1 < 0 || x1 >= tm.Width || y1 < 0 || y1 >= tm.Height ||
		x2 < 0 || x2 >= tm.Width || y2 < 0 || y2 >= tm.Height {
		return false
	}

	// We get the tile at the first position.
	tile1 := tm.GetTile(x1, y1)

	// If the tile at the first position is nil, we return false.
	if tile1 == nil {
		return false
	}

	// We get the tile at the second position.
	tile2 := tm.GetTile(x2, y2)

	// If the tile at the second position is nil, we return false.
	if tile2 == nil {
		return false
	}

	// If the tile at the first position is a wall, we return false.
	if tile1.Type == TileTypeWall {
		return false
	}

	// check every tile between the two tiles to see if they are walls or
	// closed doors. If they are, we return false.
	for _, tile := range tm.GetTilesBetween(x1, y1, x2, y2) {
		if tile.Type == TileTypeWall || tile.Type == TileTypeClosedDoor {
			return false
		}
	}

	// If we get here, we return true.
	return true
}

// GetTilesBetween returns a slice of tiles between the two given positions.
// Obviously this needs to use some cool vector math to work out what tiles are
// between the two positions. This uses the Bresenham's line algorithm to
// calculate the tiles between the two positions.
func (tm *Grid) GetTilesBetween(x1 int, y1 int, x2 int, y2 int) []Tile {
	// We create a slice of tiles to hold the tiles between the two positions.
	tiles := []Tile{}

	// We calculate the difference between the two positions.
	dx := x2 - x1
	dy := y2 - y1

	// We calculate the absolute value of the difference between the two
	// positions.
	ax := abs(dx)
	ay := abs(dy)

	// We calculate the sign of the difference between the two positions.
	sx := sign(dx)
	sy := sign(dy)

	// We calculate the error.
	err := ax - ay

	// We loop until we reach the second position.
	for {
		// We get the tile at the first position.
		tile := tm.GetTile(x1, y1)

		// If the tile is not nil, we append it to the slice of tiles.
		if tile != nil {
			tiles = append(tiles, *tile)
		}

		// If we have reached the second position, we break out of the loop.
		if x1 == x2 && y1 == y2 {
			break
		}

		// We calculate the error2.
		err2 := err * 2

		// If the error2 is greater than the negative difference between the
		// two positions, we subtract the difference from the error and
		// increment the first position by the sign of the difference between
		// the two positions.
		if err2 > -ay {
			err -= ay
			x1 += sx
		}

		// If the error2 is less than the positive difference between the two
		// positions, we add the difference to the error and increment the
		// second position by the sign of the difference between the two
		// positions.

		if err2 < ax {
			err += ax
			y1 += sy
		}
	}

	// We return the slice of tiles.
	return tiles
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

// Dump dumps an ascii representation of the grid to stdout.
//
// walls are #
// closed doors are +
// open doors are /
// floors are .
// stairs up are <
// stairs down are >
func (tm *Grid) Dump() {
	for y := 0; y < tm.Height; y++ {
		for x := 0; x < tm.Width; x++ {
			tile := tm.GetTile(x, y)
			if tile == nil {
				continue
			}
			switch tile.Type {
			case TileTypeWall:
				fmt.Printf("#")
			case TileTypeClosedDoor:
				fmt.Printf("+")
			case TileTypeOpenDoor:
				fmt.Printf("/")
			case TileTypeFloor:
				fmt.Printf(".")
			case TileTypeStairsUp:
				fmt.Printf("<")
			case TileTypeStairsDown:
				fmt.Printf(">")
			}
		}
		fmt.Println()
	}
}
