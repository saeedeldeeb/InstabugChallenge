package events

import "github.com/goravel/framework/contracts/event"

type MessageCreated struct {
}

func (receiver *MessageCreated) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
