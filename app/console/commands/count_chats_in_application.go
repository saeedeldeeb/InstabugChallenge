package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

type CountChatsInApplication struct {
}

// Signature The name and signature of the console command.
func (receiver *CountChatsInApplication) Signature() string {
	return "count:chats"
}

// Description The console command description.
func (receiver *CountChatsInApplication) Description() string {
	return "Count the number of chats in the application."
}

// Extend The console command extend.
func (receiver *CountChatsInApplication) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *CountChatsInApplication) Handle(ctx console.Context) error {
	var results []struct {
		ApplicationId uint
		ChatCount     int
	}

	err := facades.Orm().Query().
		Raw("SELECT application_id, COUNT(*) as chat_count FROM chats GROUP BY application_id").
		Scan(&results)
	if err != nil {
		facades.Log().Error(err)
		return err
	}

	for _, result := range results {
		_, err := facades.Orm().Query().Exec("UPDATE applications SET chats_count = ? WHERE id = ?", result.ChatCount, result.ApplicationId)
		if err != nil {
			facades.Log().Error(err)
			return err
		}
	}
	return nil
}
