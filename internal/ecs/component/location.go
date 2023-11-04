package component

import (
	"github.com/matjam/sword/internal/ecs"
)

// Location is the location of an entity on the Grid.
type Location struct {
	id ecs.ID

	X, Y int
}

func (*Location) New(id ecs.ID) ecs.Component {
	return &Location{id: id}
}

func (*Location) ID() ecs.ID {
	return 0
}

func (*Location) Name() string {
	return "location"
}
