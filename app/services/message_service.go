package services

import (
	"chat/app/events"
	"chat/app/models"
	"encoding/json"
	"time"

	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
)

type Message interface {
	GetMessages(appToken, chatNumber string) ([]models.Message, error)
	GetMessageByNumber(appToken, chatNumber string, messageNumber int) (models.Message, error)
	CreateMessage(appToken, chatNumber, message string) (models.Message, error)
}

type MessageService struct {
	cache cache.Driver
}

func NewMessageService() *MessageService {
	return &MessageService{
		cache: facades.Cache(),
	}
}

func (r *MessageService) GetMessages(appToken, chatNumber string) ([]models.Message, error) {
	var messages []models.Message
	remember, err := r.cache.Remember("app:"+appToken+":chat:"+chatNumber+":messages", time.Minute, func() (interface{}, error) {
		err := facades.Orm().Query().
			Where("applications.token", appToken).
			Where("chats.number", chatNumber).
			Join("JOIN chats ON messages.chat_id = chats.id").
			Join("JOIN applications ON chats.application_id = applications.id").
			With("Chat").
			Get(&messages)
		if err != nil {
			return nil, err
		}

		// Convert chats to JSON
		messagesJSON, err := json.Marshal(messages)
		if err != nil {
			return nil, err
		}
		return string(messagesJSON), nil
	})
	if err != nil {
		return nil, err
	}

	// Convert JSON back to []models.Message
	err = json.Unmarshal([]byte(remember.(string)), &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageService) GetMessageByNumber(appToken, chatNumber string, messageNumber int) (models.Message, error) {
	var message models.Message
	err := facades.Orm().Query().
		Where("applications.token", appToken).
		Where("chats.number", chatNumber).
		Where("messages.number", messageNumber).
		Join("JOIN chats ON messages.chat_id = chats.id").
		Join("JOIN applications ON chats.application_id = applications.id").
		With("Chat").
		FirstOrFail(&message)
	if err != nil {
		return models.Message{}, err
	}
	return message, nil
}

func (r *MessageService) CreateMessage(appToken, chatNumber, message string) (models.Message, error) {
	var chat models.Chat
	err := facades.Orm().Query().
		Where("applications.token", appToken).
		Where("chats.number", chatNumber).
		Join("JOIN applications ON chats.application_id = applications.id").
		FirstOrFail(&chat)
	if err != nil {
		return models.Message{}, err
	}

	// Start a transaction
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return models.Message{}, err
	}
	defer func(tx orm.Transaction) {
		if r := recover(); r != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		}
	}(tx)

	type Result struct {
		ID        int
		MaxNumber int
	}
	var result Result
	err = tx.Raw("SELECT 5 as id, COALESCE(MAX(number), 0) as MaxNumber FROM messages WHERE chat_id = ? FOR UPDATE", chat.ID).
		Scan(&result)
	if err != nil {
		return models.Message{}, err
	}

	_, err = tx.Exec("INSERT INTO messages (number, chat_id, body) VALUES (?, ?, ?)", result.MaxNumber+1, chat.ID, message)
	if err != nil {
		return models.Message{}, err
	}

	err = tx.Commit()
	if err != nil {
		return models.Message{}, err
	}

	msg, err := r.GetMessageByNumber(appToken, chatNumber, result.MaxNumber+1)
	if err != nil {
		return models.Message{}, err
	}

	// Fire an event to index the message in Elasticsearch
	messagesJSON, _ := json.Marshal(msg)
	err = facades.Event().Job(&events.MessageCreated{}, []event.Arg{
		{Type: "string", Value: string(messagesJSON)},
	}).Dispatch()
	if err != nil {
		facades.Log().Error(err)
	}

	return msg, nil
}
