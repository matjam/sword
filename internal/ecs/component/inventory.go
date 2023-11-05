package component

import "github.com/matjam/sword/internal/ecs"

type Item struct {
	Name   string
	Weight int
}

type Inventory struct {
	MaxSize     int
	MaxCapacity int

	Items []Item
}

func (*Inventory) ComponentName() ecs.ComponentName {
	return "inventory"
}
