package ecs_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
	"github.com/matjam/sword/internal/ecs/entity"
	"github.com/matjam/sword/internal/ecs/system"
)

func TestMove(t *testing.T) {
	world := ecs.NewWorld()

	// create a player entity
	player := world.AddEntity(&entity.Player{},
		&component.Location{X: 2, Y: 2},
		&component.Move{},
		&component.Render{},
		&component.Health{
			Current: 100,
			Max:     100,
		},
		&component.Inventory{},
	)

	// add a movement system
	// TODO: probably need a way to specify the order of systems
	world.AddSystem(&system.Movement{})

	// Move the player
	// TODO: maybe a registry would be good so we could refer to the components by a
	// constant value, then we could do something like:
	// movable := ecs.GetComponent[component.Move](world, player, ComponentMove)
	// though I'm not sure if that's any better than what we have now.
	movable := ecs.GetComponent[*component.Move](world, player)
	movable.X = 1
	movable.Y = 2

	// Update the world
	world.Update(1)

	// Get the player's location
	location := ecs.GetComponent[*component.Location](world, player)
	slog.Info(fmt.Sprintf("Player location: %d, %d", location.X, location.Y))

	if location.X != 3 || location.Y != 4 {
		t.Errorf("Player location should be 3, 4")
	}
}

func BenchmarkSystem(b *testing.B) {
	world := ecs.NewWorld()

	// create a player entity
	player := world.AddEntity(&entity.Player{},
		&component.Location{X: 2, Y: 2},
		&component.Move{},
		&component.Render{},
		&component.Health{
			Current: 100,
			Max:     100,
		},
		&component.Inventory{},
	)

	// add a movement system
	world.AddSystem(&system.Movement{})

	// benchmark moving the player then running update
	for n := 0; n < b.N; n++ {
		movable := ecs.GetComponent[*component.Move](world, player)
		movable.X = 1
		movable.Y = 2
		world.Update(1)
	}

	// Get the player's location
	location := world.GetComponent(player, &component.Location{}).(*component.Location)
	slog.Info(fmt.Sprintf("Player location: %d, %d", location.X, location.Y))
}
