package entity

import (
	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

// Mob is any mob entity
type Mob struct{}

func (*Mob) EntityName() ecs.EntityName {
	return "mob"
}

// New returns the player entity and its components.
func (*Mob) New() (ecs.Entity, []ecs.Component) {
	return &Mob{}, []ecs.Component{
		&component.Location{X: 5, Y: 5},
		&component.Move{},
		&component.Render{},
		&component.Damage{},
		&component.Health{
			Current: 100,
			Max:     100,
		},
		&component.Inventory{},
	}
}
