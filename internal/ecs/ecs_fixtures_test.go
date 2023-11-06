package ecs_test

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

// TestEntityWithNoComponents is an entity that has no components.

var _ ecs.Entity = &TestEntityWithNoComponents{}

type TestEntityWithNoComponents struct{}

func (*TestEntityWithNoComponents) EntityName() ecs.EntityName {
	return "test"
}

func (*TestEntityWithNoComponents) New() (ecs.Entity, []ecs.Component) {
	return &TestEntityWithNoComponents{}, []ecs.Component{}
}

// TestEntityWithComponents is an entity that has components.

var _ ecs.Entity = &TestEntityWithComponents{}

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

// TestSystemWithNoComponents is a system that has no components.

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

// TestRenderSystem is a system implements ecs.RenderSystem.

var _ = ecs.RenderSystem(&TestRenderSystem{})

type TestRenderSystem struct {
	world *ecs.World
}

func (sys *TestRenderSystem) Init(world *ecs.World) {
	sys.world = world
}

func (*TestRenderSystem) SystemName() ecs.SystemName {
	return "test_render_system"
}

func (sys *TestRenderSystem) Update(deltaTime time.Duration) {}

func (*TestRenderSystem) Components() []ecs.Component {
	return []ecs.Component{}
}

func (*TestRenderSystem) Draw(screen *ebiten.Image) {}

// TestSystemMovement is a system that implements ecs.System
// and is interested in the Move and Location components

var _ = ecs.System(&TestSystemMovement{})

type TestSystemMovement struct {
	world *ecs.World
}

// Init initializes the system.
func (sys *TestSystemMovement) Init(world *ecs.World) {
	sys.world = world
}

// SystemName returns the name of the system.
func (sys *TestSystemMovement) SystemName() ecs.SystemName {
	return "movement"
}

// Components returns the components that the system is interested in.
func (sys *TestSystemMovement) Components() []ecs.Component {
	return []ecs.Component{
		&component.Move{},
		&component.Location{},
	}
}
