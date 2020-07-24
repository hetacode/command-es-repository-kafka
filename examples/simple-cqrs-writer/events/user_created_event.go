package exampleevents

import (
	goeh "github.com/hetacode/go-eh"
)

type UserCreatedEvent struct {
	goeh.EventData
	CreateTime string `json:"create_time"`
	Version    int32  `json:"version"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (e *UserCreatedEvent) GetType() string {
	return "UserCreatedEvent"
}
