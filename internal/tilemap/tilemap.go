package tilemap

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

//go:generate go-enum --marshal

// ENUM(wall, door, corridor, floor, stairs, water, lava, trap, rubble, grass, tree, bush, rock, dirt, sand, bridge, void)
type TileType uint8

type Tile struct {
	Type TileType
}

type TileMap struct {
	Width  int
	Height int
	Tiles []Tile

	fontFace *font.Face
}

func (tm *TileMap) WithFont(fontFace *font.Face) *TileMap {
	tm.fontFace = fontFace
	return tm
}

func NewTileMap(width int, height int) *TileMap {
	tm := &TileMap{
		Width:  width,
		Height: height,
		Tiles:  make([]Tile, width*height),
	}

	for i := 0; i < width*height; i++ {
		tm.Tiles[i].Type = TileTypeWall
	}
	return tm
}

func (tm *TileMap) GetTile(x int, y int) *Tile {
	if x < 0 || x >= tm.Width || y < 0 || y >= tm.Height {
		return nil
	}
	return &tm.Tiles[y*tm.Width+x]
}

func (tm *TileMap) SetTile(x int, y int, tile *Tile) {
	if x < 0 || x >= tm.Width || y < 0 || y >= tm.Height {
		return
	}
	tm.Tiles[y*tm.Width+x] = *tile
}

func (tm *TileMap) Draw(screen *ebiten.Image) {
	for y := 0; y < tm.Height; y++ {
		for x := 0; x < tm.Width; x++ {
			tile := tm.GetTile(x, y)
			if tile == nil {
				continue
			}

			switch tile.Type {
			case TileTypeWall:
				text.Draw(screen, "█", *tm.fontFace, x*32, y*32, color.White)
			case TileTypeFloor:
				text.Draw(screen, "░", *tm.fontFace, x*32, y*32, color.White)
			}
		}
	}
}