package component

import (
	"github.com/matjam/sword/internal/ecs"
)

// Location is the location of an entity on the Grid.
type Location struct {
	X, Y int
}

func (*Location) New() ecs.Component {
	return &Location{}
}

func (*Location) Name() string {
	return "Location"
}
