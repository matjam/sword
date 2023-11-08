package tileset

import (
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matjam/sword/internal/terrain"
)

// Tileset represents a tileset atlas, for use with a tilemap and
// an autotiler. It contains the autotiles and fixtures, all of which
// are the same size and located on the same image.
type Tileset struct {
	name string
	// The image containing the tileset atlas
	atlas *ebiten.Image
	// The size of each tile in the atlas
	tileSize int
	// The number of columns in the atlas
	columns int
	// The number of rows in the atlas
	rows int
	// The autotiles in the atlas
	autotiles []*ebiten.Image
	// The fixtures in the atlas
	fixtures map[string]*ebiten.Image
}

func Load(name string,
	atlas *ebiten.Image,
	tileSize int,
	columns int, rows int,
	autotiles [][2]int,
	fixtures map[string][2]int) *Tileset {

	if len(autotiles) != 16 {
		slog.Error("autotiles must contain 16 entries", "name", name, "autotiles", len(autotiles))
	}

	ts := &Tileset{
		name:      name,
		atlas:     atlas,
		tileSize:  tileSize,
		columns:   columns,
		rows:      rows,
		autotiles: make([]*ebiten.Image, len(autotiles)),
		fixtures:  make(map[string]*ebiten.Image),
	}

	// create the autotiles
	for i, coords := range autotiles {
		x := coords[0] * tileSize
		y := coords[1] * tileSize
		ts.autotiles[i] = ts.atlas.SubImage(image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x + tileSize, Y: y + tileSize},
		}).(*ebiten.Image)
	}

	// create the fixtures
	for name, coords := range fixtures {
		x := coords[0] * tileSize
		y := coords[1] * tileSize
		ts.fixtures[name] = ts.atlas.SubImage(image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x + tileSize, Y: y + tileSize},
		}).(*ebiten.Image)
	}

	slog.Info("loaded tileset", "name", ts.name, "autotiles", len(ts.autotiles), "fixtures", len(ts.fixtures))

	return ts
}

func (ts *Tileset) Render(src *terrain.Terrain, dst *ebiten.Image, x int, y int, viewport image.Rectangle, scale int) {
	for y := 0; y < src.Height; y++ {
		for x := 0; x < src.Width; x++ {
			// don't render tiles that are outside the viewport
			if x < viewport.Min.X || x >= viewport.Max.X || y < viewport.Min.Y || y >= viewport.Max.Y {
				continue
			}

			tile := src.Get(x, y)
			if tile == terrain.Stone && !ts.isReachable(src, x, y) {
				continue
			}

			// Given the specific tile tyle (e.g. Stone, Room, Corridor, Door), render
			// the correct tile from the tileset atlas.
			//
			// We use a bitmask that represents the surrounding tiles, and use that to
			// determine which tile to render.
			//
			// the bitmask is a 4 bit number, where each bit represents a tile in one of
			// the cardinal directions. The bits are ordered like this:
			//
			//  1
			// 8 2
			//  4
			//
			// The bitmask only represents the tiles in the cardinal directions, not the
			// tile itself. For the purposes of rendering the tiles, when we render a tile
			// that is "stone", a door is considered also a solid tile so the bitmask in
			// that case would be 1 for that tile.
			//
			// The bitmask is calculated by iterating over the surrounding tiles, and
			// setting the bit in the bitmask if the tile is solid.
			//
			// For example, if the tile is surrounded by solid tiles in the north and
			// west, the bitmask would be 9 (1001).
			//
			// The bitmask is then used to index into the autotiles array, which contains
			// the correct tile to render for that bitmask.
			//
			// If the tile is not a solid tile, then we render the tile from the fixtures
			// map, which contains the correct tile to render for that tile type.
			//
			// If the tile is a solid tile but there are no surrounding solid tiles, then
			// we render the tile from the autotiles array at index 0, which is the
			// default tile for that tile type.
			//
			// Finally, if the tile is a room or corridor, we render nothing. This is
			// because we don't want to render the floor tiles for rooms and corridors,
			// as they are rendered by the room and corridor systems.

			// calculate the bitmask
			var bitmask uint8 = 0
			if tile == terrain.Stone {
				// check north
				if y > 0 && src.Get(x, y-1) == terrain.Stone && ts.isReachable(src, x, y-1) {
					bitmask |= 1
				}
				// check east
				if x < src.Width-1 && src.Get(x+1, y) == terrain.Stone && ts.isReachable(src, x+1, y) {
					bitmask |= 2
				}
				// check south
				if y < src.Height-1 && src.Get(x, y+1) == terrain.Stone && ts.isReachable(src, x, y+1) {
					bitmask |= 4
				}
				// check west
				if x > 0 && src.Get(x-1, y) == terrain.Stone && ts.isReachable(src, x-1, y) {
					bitmask |= 8
				}
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*ts.tileSize), float64(y*ts.tileSize))
			if scale != 1 {
				op.GeoM.Scale(float64(scale), float64(scale))
			}

			switch tile {
			case terrain.Stone:
				dst.DrawImage(ts.autotiles[bitmask], op)
			case terrain.Door:
				dst.DrawImage(ts.fixtures["door_unlocked"], op)
			case terrain.Room:
				dst.DrawImage(ts.fixtures["floor_dots"], op)
			case terrain.Corridor:
				dst.DrawImage(ts.fixtures["floor_checker_1"], op)
			}
		}
	}
}

func (ts *Tileset) isReachable(src *terrain.Terrain, x, y int) bool {
	// scan every tile in all 8 directions around the given tile, and if any of them
	// are not a stone tile, then the tile is reachable.

	// check north
	if y > 0 && src.Get(x, y-1) != terrain.Stone {
		return true
	}
	// check north east
	if y > 0 && x < src.Width-1 && src.Get(x+1, y-1) != terrain.Stone {
		return true
	}
	// check east
	if x < src.Width-1 && src.Get(x+1, y) != terrain.Stone {
		return true
	}
	// check south east
	if y < src.Height-1 && x < src.Width-1 && src.Get(x+1, y+1) != terrain.Stone {
		return true
	}
	// check south
	if y < src.Height-1 && src.Get(x, y+1) != terrain.Stone {
		return true
	}
	// check south west
	if y < src.Height-1 && x > 0 && src.Get(x-1, y+1) != terrain.Stone {
		return true
	}
	// check west
	if x > 0 && src.Get(x-1, y) != terrain.Stone {
		return true
	}
	// check north west
	if y > 0 && x > 0 && src.Get(x-1, y-1) != terrain.Stone {
		return true
	}

	return false
}

// all the bits in the bitmask from 0-15
//     WSEN
// 0 = 0000
// 1 = 0001
// 2 = 0010
// 3 = 0011
// 4 = 0100
// 5 = 0101
// 6 = 0110
// 7 = 0111
// 8 = 1000
// 9 = 1001
// 10 = 1010
// 11 = 1011
// 12 = 1100
// 13 = 1101
// 14 = 1110
// 15 = 1111
