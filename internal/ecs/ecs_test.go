package ecs_test

import (
	"testing"

	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

func TestECS(t *testing.T) {
	world := ecs.NewWorld()

	// register some components
	world.RegisterComponent(&component.Location{})
	world.RegisterComponent(&component.Movable{})
	world.RegisterComponent(&component.Drawable{})
	world.RegisterComponent(&component.Health{})

}
