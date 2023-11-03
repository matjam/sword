package component

import (
	"github.com/matjam/sword/internal/ecs"
)

// Health is the health of an entity.
type Health struct {
	Max     int
	Current int
}

func (*Health) New() ecs.Component {
	return &Health{}
}

func (*Health) Name() string {
	return "Health"
}

// Damage deals damage to the entity and returns the current health.
func (h *Health) Damage(d int) int {
	h.Current -= d
	if h.Current < 0 {
		h.Current = 0
	}
	return h.Current
}

// Heal heals the entity and returns the current health.
func (h *Health) Heal(d int) int {
	h.Current += d
	if h.Current > h.Max {
		h.Current = h.Max
	}
	return h.Current
}
