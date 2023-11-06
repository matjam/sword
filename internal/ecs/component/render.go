package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/matjam/sword/internal/assets"
	"github.com/matjam/sword/internal/ecs"
)

type Render struct {
	// Glyph is the rune to draw for text based rendering.
	Glyph rune
	// Color is the color to draw the glyph.
	Color color.Color
	// Sprite is the sprite to draw for sprite based rendering.
	Sprite *ebiten.Image
}

func (*Render) ComponentName() ecs.ComponentName {
	return "render"
}

// Draw draws the entity to the screen. x & y are grid coordinates.
func (d *Render) Draw(screen *ebiten.Image, x, y, gridSize int) {
	if d.Sprite != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x*gridSize), float64(y*gridSize))
		screen.DrawImage(d.Sprite, op)
	} else if d.Glyph != 0 {
		text.Draw(screen, string(d.Glyph), assets.GetFont("square"), x*gridSize, y*(gridSize-1), d.Color)
	}
}
