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

// EntityName is a unique identifier for an entity type in the ECS.
type EntityName string

// EntityID is a unique identifier for an instance of an entity in the ECS.
type EntityID ID

type ComponentName string

// ComponentID is a unique identifier for an instance of a component in the ECS.
type ComponentID ID

// SystemName is a unique identifier for an instance of a system in the ECS.
type SystemName string

// Entity is a unique object in the ECS. It is made up of a unique ID, and a
// set of components.
type Entity interface {
	// New returns a new instance of the entity, and a list of components to
	// add to the entity.
	New() (Entity, []Component)

	// EntityName returns the name of the entity type.
	EntityName() EntityName
}

// Components are the data associated with an entity. Each component can store
// data specific for a given system , and can be added to an entity to be
// processed by that system.
type Component interface {
	// ComponentName returns the name of the component.
	ComponentName() ComponentName
}

// system operates on components associated with entities.
type System interface {
	// Init is called when the system is registered with the world.
	Init(world *World)
	// SystemName returns the name of the system.
	SystemName() SystemName
	// Components returns a list of component types that this
	// system operates on.
	Components() []Component
	// Update is called every frame to update the system.
	Update(deltaTime time.Duration)
}

// World is the main ECS object. It contains all entities and systems.
//
// We need to maintain several data structures in order to efficiently query
// the world for entities that have a given set of components. We need to be
// able to query the world for entities that have a given set of components,
// as well as retrieve all components used by a given system. We also need to
// be able to retrieve a component for a given entity.
type World struct {
	// Every entity, component and system has a unique ID. The nextUniqueID
	// field stores the next ID to be used.
	nextUniqueID ID

	// entities holds all of the entities in the world. Each entity is stored
	// by its ID.
	entities map[EntityID]Entity

	// all entities of a given type
	entitiesByName map[EntityName][]EntityID

	// There can only be a single System of a given type, so we don't need a
	// registry for those.
	systems map[SystemName]System

	// components holds each instance of a component. Each component created
	// for an entity is stored here, and can be retrieved by its ID.
	components map[ComponentID]Component

	// entityComponents is a map of Entity IDs to a map of component IDs keyed
	// by component name.
	entityComponents map[EntityID]map[ComponentName]ComponentID

	// When running the main loop, a system will need to query the world for
	// all components that it operates on. We need to be able to quickly
	// retrieve all components of a given type, so we store them in a map
	// keyed bythe name of the system.
	systemComponents map[SystemName]map[ComponentName][]ComponentID

	// componentEntities is a map of component names to a list of entity IDs
	// that have that component.
	componentEntities map[ComponentName][]EntityID

	// componentGroups
}

func NewWorld() *World {
	w := &World{
		nextUniqueID:      1,
		entities:          make(map[EntityID]Entity),
		entitiesByName:    make(map[EntityName][]EntityID),
		systems:           make(map[SystemName]System),
		components:        make(map[ComponentID]Component),
		entityComponents:  make(map[EntityID]map[ComponentName]ComponentID),
		systemComponents:  make(map[SystemName]map[ComponentName][]ComponentID),
		componentEntities: make(map[ComponentName][]EntityID),
	}

	return w
}

// AddSystem adds a system to the world.
func (w *World) AddSystem(system System) {
	system.Init(w)

	w.systems[system.SystemName()] = system
	w.systemComponents[system.SystemName()] = make(map[ComponentName][]ComponentID)

	// Add the components that the system operates on to the systemComponents
	// map. When entities are added to the world, we'll add their components
	// to the systemComponents[SystemName][ComponentName] map.
	for _, component := range system.Components() {
		name := component.ComponentName()
		w.systemComponents[system.SystemName()][name] = make([]ComponentID, 0)
	}

	slog.Info("registered system", "system", system.SystemName(), "components", system.Components())
}

// AddEntity adds an entity to the world. It returns the entity ID. Optionally, you can
// pass a list of components to add to the entity.
func (w *World) AddEntity(entity Entity) EntityID {
	id := EntityID(w.nextID())

	entity, components := entity.New()

	if len(components) == 0 {
		slog.Warn("adding entity with no components", "entity", entity.EntityName())
	}

	w.entities[id] = entity
	componentNames := make([]ComponentName, 0)
	for _, component := range components {
		w.AddComponent(id, component)
		componentNames = append(componentNames, component.ComponentName())
	}

	// Add the entity to the entitiesByName map.
	if _, ok := w.entitiesByName[entity.EntityName()]; !ok {
		w.entitiesByName[entity.EntityName()] = make([]EntityID, 0)
	}
	w.entitiesByName[entity.EntityName()] = append(w.entitiesByName[entity.EntityName()], id)

	slog.Info("added entity", "id", id, "components", componentNames)
	return id
}

// AddComponent adds a component to an entity.
func (w *World) AddComponent(entityID EntityID, component Component) {
	id := ComponentID(w.nextID())
	w.components[id] = component
	name := component.ComponentName()

	// Add the component to the entity.
	if _, ok := w.entityComponents[entityID]; !ok {
		w.entityComponents[entityID] = make(map[ComponentName]ComponentID)
	}

	// check that the entity doesn't already have the component
	if _, ok := w.entityComponents[entityID][name]; ok {
		slog.Error("Entity already has component",
			"entity_id", entityID,
			"component", component.ComponentName(),
			"component_id", id)
	}

	// Add the component to the entity.
	w.entityComponents[entityID][name] = id

	// Add the component to the systemComponents map.
	for systemName, systemComponents := range w.systemComponents {
		if _, ok := systemComponents[name]; ok {
			w.systemComponents[systemName][name] = append(w.systemComponents[systemName][name], id)
		}
	}

	// Add the entity to the componentEntities map.
	w.componentEntities[name] = append(w.componentEntities[name], entityID)

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
func (w *World) HasComponents(entityID EntityID, components ...Component) bool {
	for _, component := range components {
		if !w.HasComponent(entityID, component) {
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

// EntitiesForSystem returns a list of entities that have all of the components
// that the given system operates on.
func (w *World) EntitiesForSystem(system System) []EntityID {
	return w.GetEntitiesWithComponents(system.Components()...)
}

// ComponentsForSystem returns a map of component names to a list of component
// IDs for the given system. This makes it easy to iterate over the components
// for a system.
func (w *World) ComponentsForSystem(system System) map[ComponentName][]ComponentID {
	systemName := system.SystemName()
	systemComponents := w.systemComponents[systemName]
	return systemComponents
}

// Update updates all systems in the world.
func (w *World) Update(deltaTime time.Duration) {
	for _, system := range w.systems {
		system.Update(deltaTime)
	}
}

// nextID returns the next unique ID to be used.
func (w *World) nextID() ID {
	id := w.nextUniqueID
	w.nextUniqueID++
	return id
}

// GetComponent returns the component of the given type for the given entity.
func GetComponent[T Component](world *World, entityID EntityID) T {
	var component T
	return world.GetComponent(entityID, component).(T)
}

func GetComponentID[T Component](world *World, componentID ComponentID) T {
	return world.components[componentID].(T)
}

func (world *World) GetComponentIDsForEntity(entityID EntityID) []ComponentID {
	components := make([]ComponentID, 0)
	for _, componentID := range world.entityComponents[entityID] {
		components = append(components, componentID)
	}
	return components
}

func (world *World) GetEntitiesWithComponents(components ...Component) []EntityID {
	entities := make([]EntityID, 0)
	for entityID := range world.entities {
		if world.HasComponents(entityID, components...) {
			entities = append(entities, entityID)
		}
	}
	return entities
}

// IterateComponents iterates of the components for a system, and calls the
// given function for each set of components. The function should take a map
// of component names to a component ID, one for each component that the system
// operates on.
//
// For example, if a system operates on a Move component and a Location
// component, the function will be called with a map of two components, one for
// Move and one for Location, with the ID of each component.
func (w *World) IterateComponents(system System, f func(map[ComponentName]ComponentID)) {
	systemName := system.SystemName()
	systemComponents := w.systemComponents[systemName]
	arg := make(map[ComponentName]ComponentID)

	if len(systemComponents) == 0 {
		// This is likely not an actual problem, but it's worth logging a warning
		// because you probably don't want to iterate over an empty list of
		// components. Nothing will happen.
		slog.Warn("IterateComponents called with a system that does not use components, stop that")
		return
	}

	entityCount := len(systemComponents[system.Components()[0].ComponentName()])
	for i := 0; i < entityCount; i++ {
		for componentName, componentIDs := range systemComponents {
			arg[componentName] = componentIDs[i]
		}

		f(arg)

		arg = make(map[ComponentName]ComponentID)
	}
}

func (w *World) GetEntity(entityID EntityID) Entity {
	return w.entities[entityID]
}

// GetEntity is a helper function that returns the entity of the given type
// for the given entity ID.
func GetEntity[T Entity](world *World, entityID EntityID) T {
	return world.GetEntity(entityID).(T)
}
