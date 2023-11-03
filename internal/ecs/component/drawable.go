package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/matjam/sword/internal/ecs"
)

type Drawable struct {
	// Glyph is the rune to draw for text based rendering.
	Glyph rune
	// Color is the color to draw the glyph.
	Color int
	// Sprite is the sprite to draw for sprite based rendering.
	Sprite *ebiten.Image
}

func (*Drawable) New() ecs.Component {
	return &Drawable{}
}

func (*Drawable) Name() string {
	return "Drawable"
}

// Draw draws the entity to the screen.
func (d *Drawable) Draw(screen *ebiten.Image, x, y int) {
	if d.Sprite != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(d.Sprite, op)
	} else {
		ebitenutil.DebugPrintAt(screen, string(d.Glyph), x, y+8)
	}
}
