package examplecommandhandlers

import (
	examplecommands "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/commands"
	examplerepository "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/repository"
	goeh "github.com/hetacode/go-eh"
)

type CreateUserCommandHandler struct {
	Repository *examplerepository.UsersRepository
}

func (c *CreateUserCommandHandler) Handle(event goeh.Event) {
	command := event.(*examplecommands.CreateUserCommand)
	event, err := c.Repository.Create(command.FirstName, command.LastName)
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
