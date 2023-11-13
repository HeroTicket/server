package ws

import "github.com/google/uuid"

type MessageType string

const (
	InitMessage  MessageType = "init"
	EventMessage MessageType = "event"
	ErrorMessage MessageType = "error"
)

func (t MessageType) Valid() bool {
	return t == InitMessage || t == EventMessage || t == ErrorMessage
}

type ID string

func (id ID) Valid() bool {
	_, err := uuid.Parse(string(id))
	return err == nil
}

type Message struct {
	Type  MessageType `json:"type"`
	ID    ID          `json:"id,omitempty"`
	Event Event       `json:"event,omitempty"`
}

type Status string

const (
	InProgress Status = "IN_PROGRESS"
	Error      Status = "ERROR"
	Done       Status = "DONE"
)

func (s Status) Valid() bool {
	return s == InProgress || s == Error || s == Done
}

type Event struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
	Data   any    `json:"data"`
}
