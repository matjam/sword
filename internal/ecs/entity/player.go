package entity

// Player is the player entity.
type Player struct{}

func (*Player) EntityName() string {
	return "player"
}
