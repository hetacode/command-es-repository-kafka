package examplecommands

import goeh "github.com/hetacode/go-eh"

type CreateUserCommand struct {
	goeh.EventData
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (c *CreateUserCommand) GetType() string {
	return "CreateUserCommand"
}
