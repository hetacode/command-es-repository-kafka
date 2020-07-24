package exampleevents

import (
	goeh "github.com/hetacode/go-eh"
)

type UserModifiedEvent struct {
	goeh.EventData
	CreateTime string `json:"create_time"`
	Version    int32  `json:"version"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (e *UserModifiedEvent) GetType() string {
	return "UserModifiedEvent"
}
