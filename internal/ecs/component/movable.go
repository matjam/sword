package component

import (
	"github.com/matjam/sword/internal/ecs"
)

// Movable is a component that stores the movement of an entity. An entity
// is moved by setting the X and Y values of the Movable component equal
// to the number of grid spaces to move in the X and Y directions in a
// single turn.
type Movable struct {
	id ecs.ID

	X, Y int
}

func (*Movable) New(id ecs.ID) ecs.Component {
	return &Movable{id: id}
}

func (*Movable) ID() ecs.ID {
	return 0
}

func (*Movable) Name() string {
	return "movable"
}
