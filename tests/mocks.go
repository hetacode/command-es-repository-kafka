package test

import (
	cerk "github.com/hetacode/command-es-repository-kafka"
	goeh "github.com/hetacode/go-eh"
)

// EntityCreatedEvent ...
type EntityCreatedEvent struct {
	goeh.EventData
	CreateTime string `json:"createTime"`
	Version    int32  `json:"version"`
	Message    string `json:"message"`
}

// GetType ...
func (e *EntityCreatedEvent) GetType() string {
	return "EntityCreatedEvent"
}

// MockEntity ,,,
type MockEntity struct {
	ID string

	Message string
}

// GetId ...
func (e *MockEntity) GetId() string {
	return e.ID
}

// MockProvider ...
type MockProvider struct {
	Events     []goeh.Event
	initEvents []goeh.Event
}

// SetInitEvents ...
func (p *MockProvider) SetInitEvents(events []goeh.Event) {
	p.initEvents = events
}

// FetchAllEvents ...
func (p *MockProvider) FetchAllEvents(batch int) (<-chan []goeh.Event, error) {
	c := make(chan []goeh.Event)
	go func() {
		c <- p.initEvents
		close(c)
	}()
	return c, nil
}

// SendEvents ...
func (p *MockProvider) SendEvents(events []goeh.Event) error {
	p.Events = events
	return nil
}

// Close ...
func (p *MockProvider) Close() {

}

// MockRepository ...
type MockRepository struct {
	*cerk.MemoryRepository
}

// Replay ...
func (r *MockRepository) Replay(events []goeh.Event) error {
	for _, e := range events {
		e.LoadPayload()
		switch e.GetType() {
		case "EntityCreatedEvent":
			entity := new(MockEntity)
			entity.ID = e.GetID()
			entity.Message = e.(*EntityCreatedEvent).Message
			r.AddOrModifyEntity(entity)
		}
	}

	return nil
}

// CreateFakeEvent
func (r *MockRepository) CreateFakeEvent() {
	event := &EntityCreatedEvent{
		EventData:  goeh.EventData{ID: "1"},
		Version:    1,
		Message:    "fake",
		CreateTime: "2009-11-10T23:00:00Z",
	}
	event.SavePayload(event)
	r.AddNewEvent(event)
}
