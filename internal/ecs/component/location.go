package component

import "github.com/gravestench/akara"

type Location struct {
	X, Y int
}

func (*Location) New() akara.Component {
	return &Location{}
}

var _ akara.Component = &Location{}
