package services

import (
	"chat/app/models"
	"encoding/json"
	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"time"
)

type Chat interface {
	GetChats(appToken string) ([]models.Chat, error)
	GetChatByNumber(appToken string, chatNumber int) (models.Chat, error)
	CreateChat(appToken string) (models.Chat, error)
}

type ChatService struct {
	//Dependent services
	cache cache.Driver
}

func NewChatService() *ChatService {
	return &ChatService{
		//Inject services
		cache: facades.Cache(),
	}
}

func (r *ChatService) GetChats(appToken string) ([]models.Chat, error) {
	var chats []models.Chat
	remember, err := r.cache.Remember("app:"+appToken+":chats", time.Minute, func() (interface{}, error) {
		err := facades.Orm().Query().
			Where("applications.token", appToken).
			Join("JOIN applications ON chats.application_id = applications.id").
			Get(&chats)
		if err != nil {
			return nil, err
		}

		// Convert chats to JSON
		chatsJSON, err := json.Marshal(chats)
		if err != nil {
			return nil, err
		}
		return string(chatsJSON), nil
	})
	if err != nil {
		facades.Log().Error(err)
		return nil, err
	}

	// Convert JSON back to []models.Chat
	err = json.Unmarshal([]byte(remember.(string)), &chats)
	if err != nil {
		facades.Log().Error(err)
		return nil, err
	}

	return chats, nil
}

func (r *ChatService) GetChatByNumber(appToken string, chatNumber int) (models.Chat, error) {
	var chat models.Chat
	err := facades.Orm().Query().
		Where("chats.number", chatNumber).
		Where("applications.token", appToken).
		Join("JOIN applications ON chats.application_id = applications.id").
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
