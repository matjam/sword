package ecs

import "time"

// Entity is a unique identifier for an entity in the ECS.
type Entity uint32

// Components are the data associated with an entity. Each component can store
// data specific for a given system, and can be added to an entity to be
// processed by that system.
type Component interface {
	// Name returns the name of the component.
	Name() string

	// New returns a new instance of the component.
	New() Component
}

//  system operates on components associated with entities.
type System interface {
	// Name returns the name of the system.
	Name() string

	// Update is called every frame to update the system.
	Update(deltaTime time.Duration)
}

// World is the main ECS object. It contains all entities and systems.
type World struct {
	nextEntity Entity
	entities   map[Entity][]Component
	components map[string]Component
	systems    map[string]System
}

func NewWorld() *World {
	return &World{}
}

func (w *World) AddSystem(s System) {}

// RegisterComponent registers a component with the world. This allows the
// component to be added to entities.
func (w *World) RegisterComponent(c Component) {
	if w.components == nil {
		w.components = make(map[string]Component)
	}
	w.components[c.Name()] = c
}

// RegisterSystem registers a system with the world. This allows the system to
// be updated every frame.
func (w *World) RegisterSystem(s System) {
	if w.systems == nil {
		w.systems = make(map[string]System)
	}
	w.systems[s.Name()] = s
}

// AddEntity adds an entity to the world. It returns the entity ID.
func (w *World) AddEntity() Entity {
	if w.entities == nil {
		w.entities = make(map[Entity][]Component)
	}

	id := w.nextEntity
	w.nextEntity++

	return id
}

// AddComponent adds a component to an entity.
func (w *World) AddComponent(e Entity, c Component) {
	w.entities[e] = append(w.entities[e], c)
}

// GetComponent returns the component of the given type for the given entity.
// If the entity does not have the component, it returns nil.
func (w *World) GetComponent(e Entity, c Component) Component {
	for _, component := range w.entities[e] {
		if component.Name() == c.Name() {
			return component
		}
	}
	return nil
}

// RemoveComponent removes the component of the given type from the given
// entity. If the entity does not have the component, it does nothing.
func (w *World) RemoveComponent(e Entity, c Component) {
	for i, component := range w.entities[e] {
		if component.Name() == c.Name() {
			w.entities[e] = append(w.entities[e][:i], w.entities[e][i+1:]...)
			return
		}
	}
}

// Update updates all systems in the world.
func (w *World) Update(deltaTime time.Duration) {
	for _, system := range w.systems {
		system.Update(deltaTime)
	}
}
