package services

import (
	"chat/app/models"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
)

type Chat interface {
	GetChats(appToken string) ([]models.Chat, error)
	GetChatByNumber(appToken string, chatNumber int) (models.Chat, error)
	CreateChat(appToken string) (models.Chat, error)
}

type ChatService struct {
	//Dependent services
}

func NewChatService() *ChatService {
	return &ChatService{
		//Inject services
	}
}

func (r *ChatService) GetChats(appToken string) ([]models.Chat, error) {
	var application models.Application
	err := facades.Orm().Query().Where("token", appToken).FirstOrFail(&application)
	if err != nil {
		return nil, err
	}

	var chats []models.Chat
	err = facades.Orm().Query().Where("application_id", application.ID).Get(&chats)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *ChatService) GetChatByNumber(appToken string, chatNumber int) (models.Chat, error) {
	var application models.Application
	err := facades.Orm().Query().Where("token", appToken).FirstOrFail(&application)
	if err != nil {
		return models.Chat{}, err
	}

	var chat models.Chat
	err = facades.Orm().Query().
		Where("number", chatNumber).
		Where("application_id", application.ID).
		FirstOrFail(&chat)
	if err != nil {
		return models.Chat{}, err
	}
	return chat, nil
}

func (r *ChatService) CreateChat(appToken string) (models.Chat, error) {
	var application models.Application
	err := facades.Orm().Query().Where("token", appToken).FirstOrFail(&application)
	if err != nil {
		return models.Chat{}, err
	}
	// Start transaction
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return models.Chat{}, err
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
	err = tx.Raw("SELECT 5 as id, COALESCE(MAX(number), 0) as MaxNumber FROM chats WHERE application_id = ? FOR UPDATE", application.ID).
		Scan(&result)
	facades.Log().Info(result)
	if err != nil {
		facades.Log().Error(err)
		return models.Chat{}, err
	}

	_, err = tx.Exec("INSERT INTO chats (number, application_id) VALUES (?, ?)", result.MaxNumber+1, application.ID)
	if err != nil {
		facades.Log().Error(err)
		return models.Chat{}, err
	}

	err = tx.Commit()
	if err != nil {
		return models.Chat{}, err
	}

	return r.GetChatByNumber(appToken, result.MaxNumber+1)
}
