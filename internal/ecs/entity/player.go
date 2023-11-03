package entity

import (
	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

// Player is the player entity.
type Player struct {
	// Name is the name of the player.
	Name string
	// id is the entity ID.
	id ecs.Entity
}

// NewPlayer returns a new player entity.
func NewPlayer(world *ecs.World, name string) *Player {
	p := &Player{
		Name: name,
	}

	p.id = world.AddEntity()

	// add components
	world.AddComponent(p.id, &component.Location{})

	return p
}
