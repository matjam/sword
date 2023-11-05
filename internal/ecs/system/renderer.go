package system

import (
	"time"

	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

// Ensure that we're implementing the ecs.System interface.
var _ = ecs.System(&Renderer{})

// Renderer renders all of the entities that have a Render component.
type Renderer struct {
	world *ecs.World
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
	// TODO: implement
}
