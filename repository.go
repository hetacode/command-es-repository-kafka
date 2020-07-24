package cerk

import (
	"fmt"

	goeh "github.com/hetacode/go-eh"
)

// Repository it's an abstraction for database which keeps all entities (aggregators) in theirs last state
type Repository interface {
	InitProvider(provider Provider, child Repository) error // ??
	AddOrModifyEntity(entity Entity)
	GetEntity(id string) (Entity, error)
	Replay(events []goeh.Event) error
	AddNewEvent(event goeh.Event)
	GetUncommitedChanges() []goeh.Event
	Save(events []goeh.Event) error
}

// MemoryRepository is an basic implementation of Repository which keep data in the memory. This struct waiting for inheritance by own Repository.
// Inherited implementation should contains special methods for each needed case of event type
type MemoryRepository struct {
	entities       map[string]Entity
	eventsToCommit []goeh.Event
	provider       Provider
}

// InitProvider set event store provider to repository and start restore entities
func (r *MemoryRepository) InitProvider(provider Provider, child Repository) error {
	r.provider = provider
	eventsBatch, err := provider.FetchAllEvents(20)
	if err != nil {
		return err
	}

	for events := range eventsBatch {
		child.Replay(events)
	}

	return nil
}

// AddOrModifyEntity just set new entity to collections of Entities.
// It will be replace if this id exists
func (r *MemoryRepository) AddOrModifyEntity(entity Entity) {
	r.entities[entity.GetId()] = entity
}

// GetEntity return current entity state provided by id
func (r *MemoryRepository) GetEntity(id string) (Entity, error) {
	entity := r.entities[id]
	if entity == nil {
		return nil, fmt.Errorf("Cannot find entity - id %s", id)
	}

	return entity, nil
}

// Replay method update state of entity/ies by provided events.
// This method should be override by true implementation with update cases for each event type
func (r *MemoryRepository) Replay(events []goeh.Event) error {
	return fmt.Errorf("Please implement this method in ith own way")
}

// AddNewEvent set newly created event to uncommitted list of events
func (r *MemoryRepository) AddNewEvent(event goeh.Event) {
	r.eventsToCommit = append(r.eventsToCommit, event)
}

// GetUncommitedChanges get new events which were created by changes methods
func (r *MemoryRepository) GetUncommitedChanges() []goeh.Event {
	return r.eventsToCommit
}

// Save events - so to be honest, just send events to bus provider
func (r *MemoryRepository) Save(events []goeh.Event) error {
	if err := r.provider.SendEvents(events); err != nil {
		return err
	}
	r.eventsToCommit = make([]goeh.Event, 0)

	return nil
}

// NewMemoryRepository create empty initialized instance
func NewMemoryRepository() *MemoryRepository {
	repository := new(MemoryRepository)
	repository.entities = make(map[string]Entity)
	repository.eventsToCommit = make([]goeh.Event, 0)

	return repository
}
