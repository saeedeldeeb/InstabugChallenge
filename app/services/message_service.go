package services

import (
	"chat/app/models"
	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/database/orm"
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
	var application models.Application
	err := facades.Orm().Query().Where("token", appToken).FirstOrFail(&application)
	if err != nil {
		return nil, err
	}

	var chat models.Chat
	err = facades.Orm().Query().
		Where("number", chatNumber).
		Where("application_id", application.ID).
		FirstOrFail(&chat)
	if err != nil {
		return nil, err
	}

	remember, err := r.cache.Remember("messages", 5, func() (interface{}, error) {
		var messages []models.Message
		err := facades.Orm().Query().
			Where("chat_id", chat.ID).
			With("Chat").
			Get(&messages)
		if err != nil {
			return nil, err
		}
		return messages, nil
	})
	if err != nil {
		return nil, err
	}

	return remember.([]models.Message), nil
}

func (r *MessageService) GetMessageByNumber(appToken, chatNumber string, messageNumber int) (models.Message, error) {
	var application models.Application
	err := facades.Orm().Query().Where("token", appToken).FirstOrFail(&application)
	if err != nil {
		return models.Message{}, err
	}

	var chat models.Chat
	err = facades.Orm().Query().
		Where("number", chatNumber).
		Where("application_id", application.ID).
		FirstOrFail(&chat)
	if err != nil {
		return models.Message{}, err
	}

	var message models.Message
	err = facades.Orm().Query().
		Where("number", messageNumber).
		Where("chat_id", chat.ID).
		With("Chat").
		FirstOrFail(&message)
	if err != nil {
		return models.Message{}, err
	}
	return message, nil
}

func (r *MessageService) CreateMessage(appToken, chatNumber, message string) (models.Message, error) {
	var application models.Application
	err := facades.Orm().Query().Where("token", appToken).FirstOrFail(&application)
	if err != nil {
		return models.Message{}, err
	}

	var chat models.Chat
	err = facades.Orm().Query().
		Where("number", chatNumber).
		Where("application_id", application.ID).
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

	return r.GetMessageByNumber(appToken, chatNumber, result.MaxNumber+1)
}
