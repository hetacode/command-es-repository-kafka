package examplecommandhandlers

import (
	examplecommands "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/commands"
	examplerepository "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/repository"
	goeh "github.com/hetacode/go-eh"
)

type UpdateUserCommandHandler struct {
	Repository *examplerepository.UsersRepository
}

func (c *UpdateUserCommandHandler) Handle(event goeh.Event) {
	command := event.(*examplecommands.UpdateUserCommand)
	event, err := c.Repository.Update(command.ID, command.FirstName, command.LastName)
	if err != nil {
		panic(err)
	}
	if err := c.Repository.Replay([]goeh.Event{event}); err != nil {
		panic(err)
	}
	if err := c.Repository.Save([]goeh.Event{event}); err != nil {
		panic(err)
	}
}
