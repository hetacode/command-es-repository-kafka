package examplerepository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	cerk "github.com/hetacode/command-es-repository-kafka"
	exampleevents "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/events"
	goeh "github.com/hetacode/go-eh"
)

type UsersRepository struct {
	*cerk.MemoryRepository
}

func (r *UsersRepository) Replay(events []goeh.Event) error {
	for _, e := range events {

		e.LoadPayload()
		switch e.GetType() {
		case "UserCreatedEvent":
			event := e.(*exampleevents.UserCreatedEvent)
			entity := new(UserEntity)
			entity.ID = e.GetID()
			entity.FirstName = event.FirstName
			entity.LastName = event.LastName
			r.AddOrModifyEntity(entity)
		case "UserModifiedEvent":
			event := e.(*exampleevents.UserModifiedEvent)
			entity, err := r.GetEntity(event.GetID())
			if err != nil {
				return err
			}
			userEntity := entity.(*UserEntity)
			userEntity.FirstName = event.FirstName
			userEntity.LastName = event.LastName
			r.AddOrModifyEntity(userEntity)
		}
	}

	return nil
}

func (r *UsersRepository) Create(firstName string, lastName string) (goeh.Event, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	event := &exampleevents.UserCreatedEvent{
		EventData:  goeh.EventData{ID: fmt.Sprintf("%s", id)},
		FirstName:  firstName,
		LastName:   lastName,
		CreateTime: time.Now().Format(time.RFC3339),
		Version:    0,
	}
	if err := event.SavePayload(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (r *UsersRepository) Update(id string, firstName string, lastName string) (goeh.Event, error) {
	entity, err := r.GetEntity(id)
	if err != nil {
		return nil, err
	}
	userEntity := entity.(*UserEntity)

	event := &exampleevents.UserModifiedEvent{
		EventData:  goeh.EventData{ID: fmt.Sprintf("%s", id)},
		FirstName:  IfThenElse(userEntity.FirstName != firstName, firstName, userEntity.FirstName).(string),
		LastName:   IfThenElse(userEntity.LastName != lastName, lastName, userEntity.LastName).(string),
		CreateTime: time.Now().Format(time.RFC3339),
		Version:    0,
	}
	if err := event.SavePayload(event); err != nil {
		return nil, err
	}

	return event, nil
}

func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}

	return b
}
