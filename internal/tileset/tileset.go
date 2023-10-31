package tileset

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matjam/sword/internal/tilemap"
)

// this package defines an interface for tileset renderer implementations.
// A tileset renderer is responsible for rendering a tilemap to the screen,
// which may be a simple character based renderer, or a more complex graphical
// renderer.
//
// It is the responsibility of the rendered to determine if two tiles in the
// tilemap are visible from each other, and also to determine if a tile is
// visible to the player. The tilemap does not know anything about the player
// or the camera, so it cannot make these decisions.
//
// The tileset renderer needs to keep track of 

type Renderer interface {
	// Init is called once when the renderer is first created.
	Init(*tilemap.TileMap)

	// Draw is called every frame to draw the tilemap to the screen.
	Draw(*ebiten.Image)
}