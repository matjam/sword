package system

import (
	"time"

	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

// Ensure that we're implementing the ecs.System interface.
var _ = ecs.System(&Movement{})

type Movement struct {
	world *ecs.World
}

// Init initializes the system.
func (sys *Movement) Init(world *ecs.World) {
	sys.world = world
}

// SystemName returns the name of the system.
func (sys *Movement) SystemName() ecs.SystemName {
	return "movement"
}

// Components returns the components that the system is interested in.
func (sys *Movement) Components() []ecs.Component {
	return []ecs.Component{
		&component.Move{},
		&component.Location{},
	}
}

// Update updates the system.
func (sys *Movement) Update(deltaTime time.Duration) {
	sys.world.IterateComponents(sys, func(components map[ecs.ComponentName]ecs.ComponentID) {
		location := ecs.GetComponentID[*component.Location](sys.world, components["location"])
		movable := ecs.GetComponentID[*component.Move](sys.world, components["move"])

		// move the entity
		location.X += movable.X
		location.Y += movable.Y

		// reset the movable component
		movable.X = 0
		movable.Y = 0
	})
}
