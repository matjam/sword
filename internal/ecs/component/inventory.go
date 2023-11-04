package component

import "github.com/matjam/sword/internal/ecs"

type Item struct {
	Name   string
	Weight int
}

type Inventory struct {
	id ecs.ID

	MaxSize     int
	MaxCapacity int

	Items []Item
}

func (*Inventory) New(id ecs.ID) ecs.Component {
	return &Inventory{
		Items: []Item{},
	}
}

func (i *Inventory) ID() ecs.ID {
	return i.id
}

func (*Inventory) Name() string {
	return "inventory"
}
