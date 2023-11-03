package component

import (
	"github.com/matjam/sword/internal/ecs"
)

// Movable is a component that stores the movement of an entity. An entity
// is moved by setting the X and Y values of the Movable component equal
// to the number of grid spaces to move in the X and Y directions in a
// single turn.
type Movable struct {
	X, Y int
}

func (*Movable) New() ecs.Component {
	return &Movable{}
}

func (*Movable) Name() string {
	return "Movable"
}

var _ ecs.Component = &Movable{}
