package system

import (
	"time"

	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

type Movement struct {
	world *ecs.World
}

func (*Movement) Name() string {
	return "movement"
}

func (*Movement) New(world *ecs.World) ecs.System {
	return &Movement{
		world: world,
	}
}

func (s *Movement) Components() []ecs.Component {
	return []ecs.Component{
		&component.Movable{},
		&component.Location{},
	}
}

func (s *Movement) Update(deltaTime time.Duration) {
	// get all entities with a movable and location component
	entities := s.world.EntitiesForSystem(s)

	for _, entity := range entities {
		// get the movable component
		movable := s.world.GetComponent(entity, &component.Movable{}).(*component.Movable)

		// get the location component
		location := s.world.GetComponent(entity, &component.Location{}).(*component.Location)

		// move the entity
		location.X += movable.X
		location.Y += movable.Y

		// reset the movable component
		movable.X = 0
		movable.Y = 0
	}
}
