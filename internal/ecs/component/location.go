package component

import "github.com/matjam/sword/internal/ecs"

// Location is the location of an entity on the Grid.
type Location struct {
	X, Y int
}

func (*Location) ComponentName() ecs.ComponentName {
	return "location"
}
