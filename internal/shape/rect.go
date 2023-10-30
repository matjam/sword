package shape

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
)

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func NewRect(x int, y int, width int, height int) *Rect {
	return &Rect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (r *Rect) Contains(x int, y int) bool {
	return x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}

func (r *Rect) Overlaps(other *Rect) bool {
	return r.X < other.X+other.Width && r.X+r.Width > other.X && r.Y < other.Y+other.Height && r.Y+r.Height > other.Y
}

func (r *Rect) Center() (int, int) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

func (r *Rect) CenterX() int {
	return r.X + r.Width/2
}

func (r *Rect) CenterY() int {
	return r.Y + r.Height/2
}

func (r *Rect) TopLeft() (int, int) {
	return r.X, r.Y
}

func (r *Rect) TopRight() (int, int) {
	return r.X + r.Width, r.Y
}

func (r *Rect) BottomLeft() (int, int) {
	return r.X, r.Y + r.Height
}

func (r *Rect) BottomRight() (int, int) {
	return r.X + r.Width, r.Y + r.Height
}

func (r *Rect) Top() int {
	return r.Y
}

func (r *Rect) Bottom() int {
	return r.Y + r.Height
}

func (r *Rect) Left() int {
	return r.X
}

func (r *Rect) Right() int {
	return r.X + r.Width
}

func (r *Rect) Move(x int, y int) {
	r.X += x
	r.Y += y
}

func (r *Rect) Resize(width int, height int) {
	r.Width += width
	r.Height += height
}

func (r *Rect) Clone() *Rect {
	return &Rect{
		X:      r.X,
		Y:      r.Y,
		Width:  r.Width,
		Height: r.Height,
	}
}

func (r *Rect) String() string {
	return fmt.Sprintf("Rect{X: %d, Y: %d, Width: %d, Height: %d}", r.X, r.Y, r.Width, r.Height)
}

func (r *Rect) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"X":      r.X,
		"Y":      r.Y,
		"Width":  r.Width,
		"Height": r.Height,
	})
}

func (r *Rect) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	if x, ok := v["X"]; ok {
		r.X = int(x.(float64))
	}
	if y, ok := v["Y"]; ok {
		r.Y = int(y.(float64))
	}
	if width, ok := v["Width"]; ok {
		r.Width = int(width.(float64))
	}
	if height, ok := v["Height"]; ok {
		r.Height = int(height.(float64))
	}

	return nil
}

func (r *Rect) Draw(screen *ebiten.Image, fontFace *font.Face, color color.Color) {
	for x := r.X; x < r.X+r.Width; x++ {
		for y := r.Y; y < r.Y+r.Height; y++ {
			ebitenutil.DrawRect(screen, float64(x), float64(y), 1, 1, color)
		}
	}
}