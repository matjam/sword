package tilechar

// package tilechar implements a simple character based tileset renderer.
// It is responsible for rendering a character based tilemap to the screen.

import (
	"github.com/matjam/sword/internal/tilemap"
)

type Renderer struct {
	// The tilemap to render
	tilemap *tilemap.Tilemap
}