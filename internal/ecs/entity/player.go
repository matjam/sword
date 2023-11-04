package entity

import (
	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

// Player is the player entity.
type Player struct {
	// id is the entity ID.
	id ecs.ID

	// Name is the name of the player.
	Name string
}

// NewPlayer returns a new player entity.
func NewPlayer(world *ecs.World, name string, id ecs.ID) *Player {
	p := &Player{
		id:   id,
		Name: name,
	}

	p.id = world.AddEntity()

	// add components
	world.AddComponent(p.id, &component.Location{})

	return p
}
