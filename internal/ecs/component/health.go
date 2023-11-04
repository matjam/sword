package component

import (
	"github.com/matjam/sword/internal/ecs"
)

// Health is the health of an entity.
type Health struct {
	id      ecs.ID
	Max     int
	Current int
}

func (*Health) New(id ecs.ID) ecs.Component {
	return &Drawable{id: id}
}

func (h *Health) ID() ecs.ID {
	return h.id
}

func (*Health) Name() string {
	return "health"
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
