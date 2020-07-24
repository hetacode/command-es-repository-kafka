package examplecommands

import goeh "github.com/hetacode/go-eh"

type UpdateUserCommand struct {
	goeh.EventData
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (c *UpdateUserCommand) GetType() string {
	return "UpdateUserCommand"
}
