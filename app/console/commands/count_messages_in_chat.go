package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

type CountMessagesInChat struct {
}

// Signature The name and signature of the console command.
func (receiver *CountMessagesInChat) Signature() string {
	return "count:messages"
}

// Description The console command description.
func (receiver *CountMessagesInChat) Description() string {
	return "Count the number of messages in the chat."
}

// Extend The console command extend.
func (receiver *CountMessagesInChat) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *CountMessagesInChat) Handle(ctx console.Context) error {
	var results []struct {
		ChatId       uint
		MessageCount int
	}

	err := facades.Orm().Query().
		Raw("SELECT chat_id, COUNT(*) as message_count FROM messages GROUP BY chat_id").
		Scan(&results)
	if err != nil {
		return err
	}

	for _, result := range results {
		_, err := facades.Orm().Query().Exec("UPDATE chats SET messages_count = ? WHERE id = ?", result.MessageCount, result.ChatId)
		if err != nil {
			return err
		}
	}
	return nil
}
