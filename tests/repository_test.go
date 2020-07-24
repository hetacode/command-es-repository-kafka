package test

import (
	"strings"
	"testing"
	"time"

	"github.com/go-playground/assert"
	cerk "github.com/hetacode/command-es-repository-kafka"
	goeh "github.com/hetacode/go-eh"
)

var eventsMapper *goeh.EventsMapper

func Test_Replay_Method_For_MemoryRepository_Should_Be_Override(t *testing.T) {
	repository := new(cerk.MemoryRepository)

	err := repository.Replay(make([]goeh.Event, 0))
	if err == nil {
		t.Error("Error shouldn't be null")
	}

	if !strings.HasPrefix(err.Error(), "Please implement this") {
		t.Errorf("Wrong error - %s", err.Error())
	}
}

func Test_Event_LoadPayload(t *testing.T) {
	eventsMapper := new(goeh.EventsMapper)
	eventsMapper.Register(new(EntityCreatedEvent))

	e, err := eventsMapper.Resolve(`{"message":"Just test", "createTime":"2009-11-10T23:00:00Z", "id": "1", "type": "EntityCreatedEvent"}`)

	if err != nil {
		t.Fatal(err)
	}

	event := e.(*EntityCreatedEvent)
	if event.Message != "Just test" {
		t.Fatalf("Wrong Message %s", event.Message)
	}

	date, err := time.Parse(time.RFC3339, event.CreateTime)
	if err != nil {
		t.Fatal(err)
	}

	if date.Year() != 2009 {
		t.Fatalf("Wrong CreateTime %s", event.CreateTime)
	}
}

func Test_Event_SavePayload(t *testing.T) {
	event := &EntityCreatedEvent{
		EventData:  goeh.EventData{ID: "1"},
		Version:    1,
		Message:    "It's just a test",
		CreateTime: time.Now().String(),
	}
	if err := event.SavePayload(event); err != nil {
		t.Fatal(err.Error())
	}

	if len(event.Payload) == 0 {
		t.Fatal("Payload is empty")
	}

	if !strings.Contains(event.Payload, `"createTime":`) {
		t.Fatalf("Wrong payload - %s", event.Payload)
	}
}

func Test_Replay_Function_In_MockRepository(t *testing.T) {
	event := &EntityCreatedEvent{
		EventData:  goeh.EventData{ID: "1"},
		Version:    1,
		Message:    "It's just a test",
		CreateTime: time.Now().String(),
	}

	repository := new(MockRepository)
	repository.MemoryRepository = cerk.NewMemoryRepository()

	events := []goeh.Event{event}
	if err := repository.Replay(events); err != nil {
		t.Fatal(err.Error())
	}

	entity, err := repository.GetEntity("1")
	if err != nil {
		t.Fatal(err.Error())
	}

	if entity == nil {
		t.Fatal("Cannot find entity")
	}

	if !strings.HasPrefix(entity.(*MockEntity).Message, "It's just") {
		t.Fatalf("Wrong Entity body - %s", entity.(*MockEntity).Message)
	}
}

func Test_InitProvider_And_Check_Generated_Entity(t *testing.T) {
	eventsMapper := new(goeh.EventsMapper)
	eventsMapper.Register(new(EntityCreatedEvent))

	e, err := eventsMapper.Resolve(`{"message":"Just test", "createTime":"2009-11-10T23:00:00Z", "id": "1", "type": "EntityCreatedEvent", "version": 1}`)

	event := e.(*EntityCreatedEvent)
	provider := new(MockProvider)
	provider.SetInitEvents([]goeh.Event{event})

	repository := new(MockRepository)
	repository.MemoryRepository = cerk.NewMemoryRepository()

	if err := repository.InitProvider(provider, repository); err != nil {
		t.Fatal(err.Error())
	}

	entity, err := repository.GetEntity("1")
	if err != nil {
		t.Fatal(err.Error())
	}

	if entity == nil {
		t.Fatal("Cannot find Entity")
	}

	mappedEntity := entity.(*MockEntity)
	assert.Equal(t, mappedEntity.Message, "Just test")
}

func Test_Save_Events(t *testing.T) {
	event := &EntityCreatedEvent{
		EventData: goeh.EventData{
			ID:      "1",
			Payload: `{"message":"Just test", "createTime":"2009-11-10T23:00:00Z", "id": "1", "type": "EntityCreatedEvent"}`,
		},
		Version: 1,
	}

	provider := new(MockProvider)

	repository := new(MockRepository)
	repository.MemoryRepository = cerk.NewMemoryRepository()

	if err := repository.InitProvider(provider, repository); err != nil {
		t.Fatal(err.Error())
	}

	if err := repository.Save([]goeh.Event{event}); err != nil {
		t.Fatal(err.Error())
	}

	if len(provider.Events) == 0 {
		t.Fatal("Provider events shouldn't be empty")
	}
}

func Test_AddNewEvent_In_Repository_And_Result_GetUncommitedChanges(t *testing.T) {
	provider := new(MockProvider)

	repository := new(MockRepository)
	repository.MemoryRepository = cerk.NewMemoryRepository()
	if err := repository.InitProvider(provider, repository); err != nil {
		t.Fatal(err.Error())
	}

	repository.CreateFakeEvent()
	events := repository.GetUncommitedChanges()
	if events == nil || len(events) != 1 {
		t.Fatal("Events array should exactly one element")
	}

	if err := repository.Save(events); err != nil {
		t.Fatal(err.Error())
	}

	if len(provider.Events) == 0 {
		t.Fatal("Provider events shouldn't be empty")
	}

	if len(repository.GetUncommitedChanges()) != 0 {
		t.Fatal("List of uncommitted events should be empty")
	}
}
