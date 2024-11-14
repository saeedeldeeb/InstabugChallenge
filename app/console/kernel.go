package console

import (
	"chat/app/console/commands"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/facades"
)

type Kernel struct {
}

func (kernel *Kernel) Schedule() []schedule.Event {
	return []schedule.Event{
		facades.Schedule().Command("count:chats").Hourly(),
		facades.Schedule().Command("count:messages").Hourly(),
	}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{
		&commands.CountChatsInApplication{},
		&commands.CountMessagesInChat{},
	}
}
