package grid

// package grid implements a generic grid of tiles. It can be used to
// represent a tilemap, or a grid of any other type of data.

type Grid[T any] struct {
	Width  int
	Height int

	grid []T
}

// NewGrid creates a new grid with the given width and height. The grid
// is initially filled with the zero value of the type.
func NewGrid[T any](width, height int) *Grid[T] {
	return &Grid[T]{
		Width:  width,
		Height: height,
		grid:   make([]T, width*height),
	}
}

// Get returns the value of the tile at the given position. If the position
// is outside the bounds of the grid, it returns the zero value of the type.
func (m *Grid[T]) Get(x, y int) T {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		var t T
		return t
	}

	return m.grid[y*m.Width+x]
}

// Set sets the value of the tile at the given position. If the position
// is outside the bounds of the grid, it does nothing.
func (m *Grid[T]) Set(x, y int, t T) {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return
	}

	m.grid[y*m.Width+x] = t
}

// Clear sets all the tiles in the grid to the given value. This is useful
// for clearing the grid before generating a new map.
func (m *Grid[T]) Clear(t T) {
	for i := range m.grid {
		m.grid[i] = t
	}
}

// SetRect sets all the tiles in the given rectangle to the given value.
// If the rectangle is outside the bounds of the grid, it does nothing.
func (m *Grid[T]) SetRect(x, y, w, h int, t T) {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return
	}

	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			m.Set(px, py, t)
		}
	}
}
