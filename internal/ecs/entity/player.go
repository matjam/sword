package entity

import (
	"image/color"

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
		&component.Location{},
		&component.Move{},
		&component.Render{
			Glyph: '☺',
			Color: color.RGBA{R: 64, G: 255, B: 64, A: 255},
		},
		&component.Health{
			Current: 100,
			Max:     100,
		},
		&component.Inventory{},
	}
}
