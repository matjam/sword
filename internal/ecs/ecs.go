// Package ecs implements the framework for the Entity-Component-System (ECS)
// architecture. This architecture is used to decouple data from logic, and is
// useful for games where entities can have many different types of data.
//
// The ECS architecture is made up of three main parts:
//
//  1. Entities are unique identifiers for objects in the game. They are just
//     numbers, and do not hold any data.
//  2. Components are the data associated with an entity. Each component can store
//     data specific for a given system, and can be added to an entity to be
//     processed by that system.
//  3. Systems operate on components associated with entities. They are the logic
//     of the game, and are responsible for updating the components.
//
// The World is the main ECS object. It contains all entities and systems.
// Once a World has been created, components and systems can be registered with
// it. Entities can then be created, and components added to them. Finally, the
// World can be updated every frame to update all systems.
//
// A system is registered with a list of components that it operates on. When
// the system is updated, it is passed a list of entities that have all of the
// components that it operates on. The system can then update the components
// for each entity.
//
// This implementation uses pointers to components and systems. In a ECS system
// for a real game, you'd want to store the data in a contiguous block of
// memory, and use indices instead of pointers. This would make it easier to
// iterate over the data, and would be more cache friendly. However, this
// implementation is simpler, and is good enough for now.
package ecs

import (
	"log/slog"
	"sort"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

// ID is a unique identifier for an instance of an entity, component or
// system in the ECS.
type ID uint32

// Entity is a unique object in the ECS. It is made up of a unique ID, and a
// set of components.
type Entity interface {
	// New returns a new instance of the entity.
	New(ID) Entity

	// ID returns the ID of the entity.
	ID() ID
}

// Components are the data associated with an entity. Each component can store
// data specific for a given system , and can be added to an entity to be
// processed by that system.
type Component interface {
	// New returns a new instance of the component.
	New(ID) Component

	// ID returns the ID of the component.
	ID() ID

	// Name returns the name of the component.
	Name() string
}

// system operates on components associated with entities.
type System interface {
	// New returns a new instance of the system.
	New(*World, ID) System

	// ID returns the ID of the system.
	ID() ID

	// Name returns the name of the system.
	Name() string

	// Update is called every frame to update the system.
	Update(deltaTime time.Duration)

	// Components returns a list of components that this
	// system operates on.
	Components() []Component
}

// World is the main ECS object. It contains all entities and systems.
type World struct {
	nextID ID

	// entities holds a list of components for each entity. The value stored
	// is a list of component IDs.
	entities map[ID]Entity

	// registries hold a map of names to IDs, and IDs to components and systems.
	// This is used to quickly find the ID of a component or system given its
	// name, and to find the component or system given its ID.
	componentRegistry     map[string]ID
	componentRegistryByID map[ID]Component
	systemRegistry        map[string]ID
	systemRegistryByID    map[ID]System

	// entityComponents is a map of component ID to a list of entity IDs.
	// This is used to quickly find all entities that have a given component.
	entityComponents map[ID]mapset.Set[ID]
}

func NewWorld() *World {
	return &World{
		nextID:           0,
		entities:         make(map[ID]Entity),
		components:       make(map[ID]Component),
		componentIDs:     make(map[string]ID),
		systems:          make(map[ID]System),
		systemIDs:        make(map[string]ID),
		systemComponents: make(map[ID][]ID),
		entityComponents: make(map[ID][]ID),
	}
}

// RegisterComponent registers a component with the world. This allows the
// component to be added to entities.
func (w *World) RegisterComponent(c Component, newComponentFunc NewComponentFunc) {
	id := w.nextID
	w.nextID++

	w.components[id] = c
	w.componentIDs[c.Name()] = id

	slog.Info("Registered", "id", id, "component", c.Name())
}

// RegisterSystem registers a system with the world. This allows the system to
// be updated every frame.
func (w *World) RegisterSystem(s System) {
	system := s.New(w) // There can be only one
	id := w.nextID
	w.nextID++

	w.systems[id] = system
	w.systemIDs[system.Name()] = id
	w.systemComponents[id] = make([]ID, 0)

	// generate the list of component IDs for this system
	for _, component := range system.Components() {
		w.systemComponents[id] = append(w.systemComponents[id], w.componentIDs[component.Name()])
	}

	slog.Info("Registered", "id", id, "system", system.Name(), "components", w.systemComponents[id])
}

// AddEntity adds an entity to the world. It returns the entity ID.
func (w *World) AddEntity() ID {
	id := w.nextID
	w.nextID++

	if w.entityComponents == nil {
		w.entityComponents = make(map[ID][]ID)
	}

	slog.Info("Added entity", "id", id)
	return id
}

// AddComponent adds a component to an entity.
func (w *World) AddComponent(e ID, c Component) {
	component := c.New()

	w.entities[e] = append(w.entities[e], c)
	w.entityComponents[c.Name()] = append(w.entityComponents[c.Name()], e)

	// check every system to see if the entity now has all of the components
	// that the system operates on. If it does, we add it to the list of
	// entities for that system.

	for _, componentNames := range w.systemComponents {
		// split the component names into a list
		components := strings.Split(componentNames, ",")
		// sort the list
		sortedComponents := sort.StringSlice(components)
		sortedComponents.Sort()

		hasAllComponents := true
		for _, componentName := range sortedComponents {
			if !w.HasComponent(e, w.components[componentName]) {
				hasAllComponents = false
				break
			}
		}

		if hasAllComponents {
			w.entityComponents[componentNames] = append(w.entityComponents[componentNames], e)
		}
	}

	slog.Info("Added component", "entity", e, "component", c.Name())
}

// HasComponent returns true if the given entity has the given component.
func (w *World) HasComponent(e ID, c Component) bool {
	for _, component := range w.entities[e] {
		if component.Name() == c.Name() {
			return true
		}
	}
	return false
}

// GetComponent returns the component of the given type for the given entity.
// If the entity does not have the component, it returns nil.
func (w *World) GetComponent(e ID, c Component) Component {
	for _, component := range w.entities[e] {
		if component.Name() == c.Name() {
			return component
		}
	}
	return nil
}

// RemoveComponent removes the component of the given type from the given
// entity. If the entity does not have the component, it does nothing.
func (w *World) RemoveComponent(e ID, c Component) {
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

// EntitiesForSystem returns a list of entities that have all of the components
// that the given system operates on.
func (w *World) EntitiesForSystem(s System) []ID {
	components := w.systemComponents[s.Name()]
	entities := w.entityComponents[components]
	return entities
}
