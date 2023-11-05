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
	world := ecs.NewWorld(
		&component.Location{},
		&component.Move{},
		&component.Render{},
		&component.Health{},
		&component.Inventory{},
	)

	// add a movement system
	// TODO: probably need a way to specify the order of systems
	world.AddSystem(&system.Movement{})

	// create a player entity
	player := world.AddEntity(&entity.Player{})
	mob := world.AddEntity(&entity.Mob{})

	// Move the player
	movable := ecs.GetComponent[*component.Move](world, player)
	movable.X = 1
	movable.Y = 2

	// Move the mob
	movable = ecs.GetComponent[*component.Move](world, mob)
	movable.X = 3
	movable.Y = 4

	// Update the world
	world.Update(1)

	ecs.Spew(world)

	// Get the player's location
	playerLocation := ecs.GetComponent[*component.Location](world, player)
	slog.Info(fmt.Sprintf("Player location: %d, %d", playerLocation.X, playerLocation.Y))

	if playerLocation.X != 3 || playerLocation.Y != 4 {
		t.Errorf("Player location should be 3, 4")
	}

	// Get the mob's location
	mobLocation := ecs.GetComponent[*component.Location](world, mob)
	slog.Info(fmt.Sprintf("Mob location: %d, %d", mobLocation.X, mobLocation.Y))

	if mobLocation.X != 8 || mobLocation.Y != 9 {
		t.Errorf("Mob location should be 8, 9")
	}
}

type TestEntityWithNoComponents struct{}

func (*TestEntityWithNoComponents) EntityName() ecs.EntityName {
	return "test"
}

func (*TestEntityWithNoComponents) New() (ecs.Entity, []ecs.Component) {
	return &TestEntityWithNoComponents{}, []ecs.Component{}
}

func TestAddEntityWithNoComponents(t *testing.T) {
	world := ecs.NewWorld()
	testEntityID := world.AddEntity(&TestEntityWithNoComponents{})

	// Get the testEntity's components
	components := world.GetComponentIDsForEntity(testEntityID)

	if len(components) != 0 {
		t.Errorf("Player should have no components")
	}
}

func TestAddDuplicateComponents(t *testing.T) {
	world := ecs.NewWorld()
	testEntityID := world.AddEntity(&TestEntityWithNoComponents{})

	// Add a duplicate component
	world.AddComponent(testEntityID, &component.Location{X: 1, Y: 1})

	t.Run("Runtime error expected", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Errorf("The code did not panic")
			}
		}()

		world.AddComponent(testEntityID, &component.Location{X: 1, Y: 1})
	})
}
