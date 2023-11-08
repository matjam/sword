package mapgen

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/matjam/sword/internal/terrain"
)

////////////////////////////////////////////////////////////////////////////////
// Drawing

func (mg *MapGenerator) DrawDebug(screen *ebiten.Image) {
	for y := 0; y < mg.Height; y++ {
		for x := 0; x < mg.Width; x++ {
			t := mg.terrainGrid.Get(x, y)
			r := mg.regionGrid.Get(x, y)

			clr := color.Color(color.RGBA{0x50, 0x50, 0x50, 0xff})
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
				mg.drawTile(screen, x, y, color.RGBA{0x70, 0x30, 0x30, 0xff})
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
