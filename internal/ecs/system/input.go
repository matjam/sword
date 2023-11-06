package system

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
)

// Ensure that we're implementing the ecs.System interface.
var _ = ecs.System(&Input{})

type Input struct {
	world  *ecs.World
	Player ecs.EntityID
	keys   []ebiten.Key
}

// Init initializes the system.
func (sys *Input) Init(world *ecs.World) {
	sys.world = world
	sys.keys = make([]ebiten.Key, 0, 20)
}

// SystemName returns the name of the system.
func (sys *Input) SystemName() ecs.SystemName {
	return "input"
}

// Components returns the components that the system is interested in.
func (sys *Input) Components() []ecs.Component {
	return []ecs.Component{
		&component.Location{},
	}
}

// Update updates the system.
func (sys *Input) Update(deltaTime time.Duration) {
	sys.keys = inpututil.AppendPressedKeys(sys.keys[:0])
	for _, key := range sys.keys {
		switch key {
		case ebiten.KeyW:
			if inpututil.IsKeyJustPressed(ebiten.KeyW) {
				sys.movePlayer(0, -1)
			}
		case ebiten.KeyS:
			if inpututil.IsKeyJustPressed(ebiten.KeyS) {
				sys.movePlayer(0, 1)
			}
		case ebiten.KeyA:
			if inpututil.IsKeyJustPressed(ebiten.KeyA) {
				sys.movePlayer(-1, 0)
			}
		case ebiten.KeyD:
			if inpututil.IsKeyJustPressed(ebiten.KeyD) {
				sys.movePlayer(1, 0)
			}
		}
	}
}

func (sys *Input) movePlayer(x, y int) {
	movable := ecs.GetComponent[*component.Move](sys.world, sys.Player)
	movable.X = x
	movable.Y = y
}
