package system

import (
	"time"

	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

type Movement struct{}

func (*Movement) SystemName() ecs.SystemName {
	return "movement"
}

func (s *Movement) Components() []ecs.Component {
	return []ecs.Component{
		&component.Move{},
		&component.Location{},
	}
}

func (s *Movement) Update(world *ecs.World, deltaTime time.Duration) {
	world.IterateComponents(s, func(components map[ecs.ComponentName]ecs.ComponentID) {
		location := ecs.GetComponentID[*component.Location](world, components["location"])
		movable := ecs.GetComponentID[*component.Move](world, components["move"])

		// move the entity
		location.X += movable.X
		location.Y += movable.Y

		// reset the movable component
		movable.X = 0
		movable.Y = 0
	})
}
