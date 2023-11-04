package component

type Item struct {
	Name   string
	Weight int
}

type Inventory struct {
	MaxSize     int
	MaxCapacity int

	Items []Item
}

func (*Inventory) ComponentName() string {
	return "inventory"
}
