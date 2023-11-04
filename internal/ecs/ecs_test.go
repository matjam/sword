package ecs_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
	"github.com/matjam/sword/internal/ecs/system"
)

func TestMovable(t *testing.T) {
	world := ecs.NewWorld()

	// register some components
	world.RegisterComponent(&component.Location{})
	world.RegisterComponent(&component.Movable{})
	world.RegisterComponent(&component.Drawable{})
	world.RegisterComponent(&component.Health{})
	world.RegisterComponent(&component.Inventory{})

	// register some systems
	world.RegisterSystem(&system.Movement{})

	// create a player entity
	player := world.AddEntity()
	world.AddComponent(player, &component.Location{})
	world.AddComponent(player, &component.Movable{})
	world.AddComponent(player, &component.Drawable{})
	world.AddComponent(player, &component.Health{})
	world.AddComponent(player, &component.Inventory{})

	// Move the player
	movable := world.GetComponent(player, &component.Movable{}).(*component.Movable)
	movable.X = 1
	movable.Y = 2

	// Update the world
	world.Update(1)

	// Get the player's location
	location := world.GetComponent(player, &component.Location{}).(*component.Location)
	slog.Info(fmt.Sprintf("Player location: %d, %d", location.X, location.Y))

	if location.X != 1 || location.Y != 1 {
		t.Errorf("Player location should be 1, 2")
	}
}
