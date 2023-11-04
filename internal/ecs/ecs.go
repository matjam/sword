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
	"time"
)

// These IDs are globally unique identifiers for entities, components and
// systems. They are used to identify an entity, component or system when
// registering them with the world, and when adding them to an entity.
type ID uint32

// EntityFactoryID is a unique identifier for an entity type in the ECS.
type EntityFactoryID ID

// EntityID is a unique identifier for an instance of an entity in the ECS.
type EntityID ID

// ComponentFactoryID is a unique identifier for a component type in the ECS.
type ComponentFactoryID ID

// ComponentID is a unique identifier for an instance of a component in the ECS.
type ComponentID ID

// SystemID is a unique identifier for an instance of a system in the ECS.
type SystemID ID

// Entity is a unique object in the ECS. It is made up of a unique ID, and a
// set of components.
type Entity interface {
	// EntityName returns the name of the entity type.
	EntityName() string
}

// Components are the data associated with an entity. Each component can store
// data specific for a given system , and can be added to an entity to be
// processed by that system.
type Component interface {
	// ComponentName returns the name of the component.
	ComponentName() string
}

// system operates on components associated with entities.
type System interface {
	// SystemName returns the name of the system.
	SystemName() string

	// Update is called every frame to update the system.
	Update(world *World, deltaTime time.Duration)

	// Components returns a list of component types that this
	// system operates on.
	Components() []Component
}

// World is the main ECS object. It contains all entities and systems.
type World struct {
	// Every entity, component and system has a unique ID. The nextUniqueID
	// field stores the next ID to be used.
	nextUniqueID ID

	// entities holds a list of components for each entity. The value stored
	// is a list of component IDs.
	entities map[EntityID]Entity

	// There can only be a single System of a given type, so we don't need a
	// Factory for those.
	systems map[SystemID]System

	// components holds each instance of a component. Each component created
	// for an entity is stored here, and can be retrieved by its ID.
	components map[ComponentID]Component

	// entityComponents is a map of Entity IDs to a map of component IDs keyed
	// by component name.
	entityComponents map[EntityID]map[string]ComponentID

	// componentEntities is a map of Component names to an array of Entity IDs.
	componentEntities map[string][]EntityID
}

func NewWorld() *World {
	return &World{
		nextUniqueID:      0,
		entities:          make(map[EntityID]Entity),
		systems:           make(map[SystemID]System),
		components:        make(map[ComponentID]Component),
		entityComponents:  make(map[EntityID]map[string]ComponentID),
		componentEntities: make(map[string][]EntityID),
	}
}

func (w *World) AddSystem(system System) {
	id := SystemID(w.nextID())

	w.systems[id] = system

	slog.Info("Registered", "id", id, "system", system.SystemName(), "components", system.Components())
}

// AddEntity adds an entity to the world. It returns the entity ID. Optionally, you can
// pass a list of components to add to the entity.
func (w *World) AddEntity(entity Entity, components ...Component) EntityID {
	id := EntityID(w.nextID())

	w.entities[id] = entity

	for _, component := range components {
		w.AddComponent(id, component)
	}

	slog.Info("Added entity", "id", id, "components", components)
	return id
}

// AddComponent adds a component to an entity.
func (w *World) AddComponent(entityID EntityID, component Component) {
	id := ComponentID(w.nextID())

	w.components[id] = component

	if _, ok := w.entityComponents[entityID]; !ok {
		w.entityComponents[entityID] = make(map[string]ComponentID)
	}
	w.entityComponents[entityID][component.ComponentName()] = id

	if _, ok := w.componentEntities[component.ComponentName()]; !ok {
		w.componentEntities[component.ComponentName()] = make([]EntityID, 0)
	}
	w.componentEntities[component.ComponentName()] =
		append(w.componentEntities[component.ComponentName()], entityID)

	slog.Info("Added component",
		"entity_id", entityID,
		"component", component.ComponentName(),
		"component_id", id)
}

// HasComponent returns true if the given entity has the given component.
func (w *World) HasComponent(entityID EntityID, component Component) bool {
	name := component.ComponentName()
	if _, ok := w.entityComponents[entityID]; ok {
		if _, ok := w.entityComponents[entityID][name]; ok {
			return true
		}
	}

	return false
}

// HasComponents returns true if the given entity has all of the given
// components.
func (w *World) HasComponents(entityID EntityID, components []string) bool {
	for _, name := range components {
		if _, ok := w.entityComponents[entityID][name]; !ok {
			return false
		}
	}

	return true
}

// GetComponent returns the component of the given type for the given entity.
// If the entity does not have the component, it returns nil.
func (w *World) GetComponent(entityID EntityID, component Component) Component {
	name := component.ComponentName()
	if _, ok := w.entityComponents[entityID]; ok {
		if componentID, ok := w.entityComponents[entityID][name]; ok {
			return w.components[componentID]
		}
	}

	return nil
}

// Update updates all systems in the world.
func (w *World) Update(deltaTime time.Duration) {
	for _, system := range w.systems {
		system.Update(w, deltaTime)
	}
}

// EntitiesForSystem returns a list of entities that have all of the components
// that the given system operates on.
func (w *World) EntitiesForSystem(system System) []EntityID {
	entities := make([]EntityID, 0)

	// Get the list of component names that the system operates on.
	componentNames := make([]string, 0)
	for _, component := range system.Components() {
		componentNames = append(componentNames, component.ComponentName())
	}

	// Get the list of entities that have all of the components that the system
	// operates on.
	for _, entityID := range w.componentEntities[componentNames[0]] {
		// If the entity does not have any of the components, we skip it.
		if !w.HasComponents(entityID, componentNames) {
			continue
		}

		entities = append(entities, entityID)
	}

	return entities
}

func (w *World) nextID() ID {
	id := w.nextUniqueID
	w.nextUniqueID++
	return id
}

func GetComponent[T Component](world *World, entityID EntityID, component Component) T {
	return world.GetComponent(entityID, component).(T)
}
