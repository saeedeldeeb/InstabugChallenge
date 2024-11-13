package controllers

import (
	"chat/app/http/transformers"
	"chat/app/models"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ChatController struct {
	//Dependent services
}

func NewChatController() *ChatController {
	return &ChatController{
		//Inject services
	}
}

func (r *ChatController) Index(ctx http.Context) http.Response {
	var chats []models.Chat
	err := facades.Orm().Query().Get(&chats)
	if err != nil {
		return nil
	}
	return ctx.Response().Json(http.StatusOK, transformers.ChatsCollectionResponse(chats))
}

func (r *ChatController) Show(ctx http.Context) http.Response {
	var application models.Application
	err := facades.Orm().Query().Where("token", ctx.Request().Input("token")).FirstOrFail(&application)

	var chat models.Chat
	err = facades.Orm().Query().
		Where("number", ctx.Request().Input("number")).
		Where("application_id", application.ID).
		FirstOrFail(&chat)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.ChatResponse(chat))
}

func (r *ChatController) Store(ctx http.Context) http.Response {
	var application models.Application
	err := facades.Orm().Query().Where("token", ctx.Request().Input("token")).FirstOrFail(&application)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	_, err = facades.Orm().Query().Exec("INSERT INTO chats (number, application_id) SELECT COALESCE(MAX(number), 0)+1, ? FROM chats WHERE application_id = ?", application.ID, application.ID)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}
	return ctx.Response().Redirect(http.StatusCreated, "/api/applications/"+ctx.Request().Input("token")+"/chats")
}
