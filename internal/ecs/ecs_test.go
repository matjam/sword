package ecs_test

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
	"github.com/matjam/sword/internal/ecs/entity"
	"github.com/matjam/sword/internal/ecs/system"
)

type TestEntityWithNoComponents struct{}

func (*TestEntityWithNoComponents) EntityName() ecs.EntityName {
	return "test"
}

func (*TestEntityWithNoComponents) New() (ecs.Entity, []ecs.Component) {
	return &TestEntityWithNoComponents{}, []ecs.Component{}
}

type TestEntityWithComponents struct{}

func (*TestEntityWithComponents) EntityName() ecs.EntityName {
	return "test"
}

func (*TestEntityWithComponents) New() (ecs.Entity, []ecs.Component) {
	return &TestEntityWithComponents{}, []ecs.Component{
		&component.Location{X: 1, Y: 1},
		&component.Move{X: 1, Y: 1},
		&component.Render{},
		&component.Health{Current: 100, Max: 100},
	}
}

func TestMove(t *testing.T) {
	world := ecs.NewWorld()

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
	world.AddComponent(testEntityID, &component.Location{X: 1, Y: 1})
	components := world.GetComponentIDsForEntity(testEntityID)

	if len(components) != 1 {
		t.Errorf("Player should have 1 component")
	}

}

func TestWorld_HasComponent(t *testing.T) {
	// Test that the HasComponent function works

	world := ecs.NewWorld()
	testEntityID := world.AddEntity(&TestEntityWithNoComponents{})

	// Add a component
	world.AddComponent(testEntityID, &component.Location{X: 1, Y: 1})

	// Test that the component exists
	if !world.HasComponent(testEntityID, &component.Location{}) {
		t.Errorf("The component should exist")
	}

	// Test that a non-existent component does not exist
	if world.HasComponent(testEntityID, &component.Move{}) {
		t.Errorf("The component should not exist")
	}
}

func TestWorld_HasComponents(t *testing.T) {
	// test that the HasComponents function works by adding multiple components
	// and checking that they exist

	world := ecs.NewWorld()
	testEntityID := world.AddEntity(&TestEntityWithComponents{})

	// Test that the components exist
	if !world.HasComponents(testEntityID, &component.Location{}, &component.Move{}, &component.Render{}, &component.Health{}) {
		t.Errorf("The components should exist")
	}

	// Test that a non-existent component does not exist
	if world.HasComponents(testEntityID, &component.Location{}, &component.Move{}, &component.Render{}, &component.Health{}, &component.Inventory{}) {
		t.Errorf("The components should not exist")
	}

	// Test that checking for a smaller set of components works
	if !world.HasComponents(testEntityID, &component.Location{}, &component.Move{}, &component.Render{}) {
		t.Errorf("The components should exist")
	}
}

func TestWorld_GetComponent(t *testing.T) {
	// Test that the GetComponent function works

	world := ecs.NewWorld()
	testEntityID := world.AddEntity(&TestEntityWithComponents{})

	// Test that the component exists
	location := ecs.GetComponent[*component.Location](world, testEntityID)

	if location.X != 1 || location.Y != 1 {
		t.Errorf("The component should exist")
	}

	t.Run("Runtime error expected", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Errorf("The code did not panic")
			}
		}()

		// Test that a non-existent component does not exist
		inventory := ecs.GetComponent[*component.Inventory](world, testEntityID)
		if inventory != nil {
			t.Errorf("The component should not exist")
		}
	})
}

func TestWorld_EntitiesForSystem(t *testing.T) {
	// Test that the EntitiesForSystem function works

	world := ecs.NewWorld()
	testEntityID := world.AddEntity(&TestEntityWithComponents{})

	// Test that the component exists
	entities := world.EntitiesForSystem(&system.Movement{})

	if len(entities) != 1 {
		t.Fatal("There should be 1 entity")
	}

	if entities[0] != testEntityID {
		t.Errorf("The entity ID should match")
	}
}

func TestWorld_ComponentsForSystem(t *testing.T) {
	// Test that the ComponentsForSystem function works

	world := ecs.NewWorld()
	world.AddSystem(&system.Movement{})
	world.AddEntity(&TestEntityWithComponents{})

	// Test that the component exists
	components := world.ComponentsForSystem(&system.Movement{})

	// Should be two components returned: Location and Move
	if len(components) != 2 {
		t.Fatal("There should be 2 component")
	}

	location := ecs.GetComponentID[*component.Location](world, components["location"][0])
	move := ecs.GetComponentID[*component.Move](world, components["move"][0])

	if location == nil || move == nil {
		t.Fatal("Location and Move components should exist")
	}

	if location.X != 1 || location.Y != 1 {
		t.Errorf("Location should be 1, 1")
	}

	if move.X != 1 || move.Y != 1 {
		t.Errorf("Move should be 1, 1")
	}
}

func TestWorld_GetComponentIDsForEntity(t *testing.T) {
	// Test that the GetComponentIDsForEntity function works

	world := ecs.NewWorld()
	testEntityID := world.AddEntity(&TestEntityWithComponents{})

	// Test that the component exists
	components := world.GetComponentIDsForEntity(testEntityID)

	if len(components) != 4 {
		t.Fatal("There should be 4 components")
	}
}

var _ ecs.System = &TestSystemWithNoComponents{}

type TestSystemWithNoComponents struct {
	world *ecs.World
}

func (sys *TestSystemWithNoComponents) Init(world *ecs.World) {
	sys.world = world
}

func (*TestSystemWithNoComponents) SystemName() ecs.SystemName {
	return "test"
}

func (sys *TestSystemWithNoComponents) Update(deltaTime time.Duration) {
	sys.world.IterateComponents(sys, func(components map[ecs.ComponentName]ecs.ComponentID) {
		// do nothing
	})
}

func (*TestSystemWithNoComponents) Components() []ecs.Component {
	return []ecs.Component{}
}

func TestWorld_AddSystemWithNoComponents(t *testing.T) {
	world := ecs.NewWorld()
	world.AddSystem(&TestSystemWithNoComponents{})
	world.Update(1)
}

func TestWorld_GetEntity(t *testing.T) {
	// Test that the GetEntity function works

	world := ecs.NewWorld()
	testEntityID := world.AddEntity(&TestEntityWithComponents{})

	// Test that the component exists
	entity := world.GetEntity(testEntityID)

	if entity == nil {
		t.Fatal("The entity should exist")
	}
}

func TestGetEntity(t *testing.T) {
	// Test that the GetEntity function works

	world := ecs.NewWorld()
	testEntityID := world.AddEntity(&TestEntityWithComponents{})

	// Test that the component exists
	entity := ecs.GetEntity[*TestEntityWithComponents](world, testEntityID)

	if entity == nil {
		t.Fatal("The entity should exist")
	}
}
