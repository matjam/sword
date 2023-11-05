package entity

import (
	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

// Player is the player entity.
type Player struct{}

func (*Player) EntityName() ecs.EntityName {
	return "player"
}

// New returns the player entity and its components.
func (*Player) New() (ecs.Entity, []ecs.Component) {
	return &Player{}, []ecs.Component{
		&component.Location{X: 2, Y: 2},
		&component.Move{},
		&component.Render{},
		&component.Health{
			Current: 100,
			Max:     100,
		},
		&component.Inventory{},
	}
}
