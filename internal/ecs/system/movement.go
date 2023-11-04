package system

import (
	"time"

	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

type Movement struct{}

func (*Movement) SystemName() string {
	return "movement"
}

func (s *Movement) Components() []ecs.Component {
	return []ecs.Component{
		&component.Move{},
		&component.Location{},
	}
}

func (s *Movement) Update(world *ecs.World, deltaTime time.Duration) {
	// get all entities with a movable and location component
	entities := world.EntitiesForSystem(s)

	for _, entity := range entities {
		// get the movable component
		movable := world.GetComponent(entity, &component.Move{}).(*component.Move)

		// get the location component
		location := world.GetComponent(entity, &component.Location{}).(*component.Location)

		// move the entity
		location.X += movable.X
		location.Y += movable.Y

		// reset the movable component
		movable.X = 0
		movable.Y = 0
	}
}
