package system

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

// Ensure that we're implementing the ecs.System interface.
var _ = ecs.RenderSystem(&Renderer{})

// Renderer renders all of the entities that have a Render component.
type Renderer struct {
	world *ecs.World

	GridSize int
}

// Init initializes the system.
func (sys *Renderer) Init(world *ecs.World) {
	sys.world = world
}

// SystemName returns the name of the system.
func (sys *Renderer) SystemName() ecs.SystemName {
	return "renderer"
}

// Components returns the components that the system is interested in.
func (sys *Renderer) Components() []ecs.Component {
	return []ecs.Component{
		&component.Render{},
		&component.Location{},
	}
}

// Update updates the system.
func (sys *Renderer) Update(delta time.Duration) {
	// the renderer system doesn't need to update anything
}

func (sys *Renderer) WillDraw() bool {
	return true
}

func (sys *Renderer) Draw(screen *ebiten.Image) {
	sys.world.IterateComponents(sys, func(components map[ecs.ComponentName]ecs.ComponentID) {
		render := ecs.GetComponentID[*component.Render](sys.world, components["render"])
		location := ecs.GetComponentID[*component.Location](sys.world, components["location"])

		render.Draw(screen, location.X, location.Y, sys.GridSize)
	})
}
