package component

// Location is the location of an entity on the Grid.
type Location struct {
	X, Y int
}

func (*Location) ComponentName() string {
	return "location"
}
