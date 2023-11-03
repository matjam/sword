package text

// package text implements a simple text based tileset renderer. It will render
// a given Grid using the font given to it.

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/matjam/sword/internal/assets"
	"github.com/matjam/sword/internal/tilemap"
	"golang.org/x/image/font"
)

type Renderer struct {
	// The tilemap to render
	tilemap *tilemap.Grid
	// The font to use for rendering
	tilefont font.Face
	// The size of the font
	size int
}

func NewRenderer(tilemap *tilemap.Grid, fontName string) tilemap.Renderer {
	return &Renderer{
		tilemap:  tilemap,
		tilefont: assets.GetFont(fontName),
		size:     assets.GetFontSize(fontName),
	}
}

// Draw the tilemap to the given destination image. The viewport is the
// rectangle of the tilemap to render.
func (r *Renderer) Draw(dst *ebiten.Image, x int, y int, viewport tilemap.Rectangle) {
	// Iterate over the tiles in the viewport, and write them to the destination,
	// line by line. We use the tilemap's width to calculate the position of the
	// tile in the tilemap.

	row := make([]rune, viewport.Width)
	destY := y

	for y := viewport.Y; y < viewport.Y+viewport.Height; y++ {
		for x := viewport.X; x < viewport.X+viewport.Width; x++ {
			tile := r.tilemap.GetTile(x, y)
			if tile == nil {
				continue
			}

			row[x-viewport.X] = tileTypeToRune[tile.Type]
		}
		text.Draw(dst, string(row), r.tilefont, x, destY, color.White)
		destY += r.size - 1

		// it doesn't matter if we don't clear the row, because we're going to
		// overwrite it anyway.
	}
}

var tileTypeToRune = map[tilemap.TileType]rune{
	tilemap.TileTypeWall:       '█',
	tilemap.TileTypeClosedDoor: '▒',
	tilemap.TileTypeOpenDoor:   '░',
	tilemap.TileTypeFloor:      ' ',
	tilemap.TileTypeStairsUp:   '<',
	tilemap.TileTypeStairsDown: '>',
}
