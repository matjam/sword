package component

import "github.com/matjam/sword/internal/ecs"

// Move is a component that stores the movement of an entity. An entity
// is moved by setting the X and Y values of the Move component equal
// to the number of grid spaces to move in the X and Y directions in a
// single turn.
type Move struct {
	X, Y int
}

func (*Move) ComponentName() ecs.ComponentName {
	return "move"
}
