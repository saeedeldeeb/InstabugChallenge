package controllers

import (
	"chat/app/models"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type MessageController struct {
	//Dependent services
}

func NewMessageController() *MessageController {
	return &MessageController{
		//Inject services
	}
}

func (r *MessageController) Index(ctx http.Context) http.Response {
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

	var messages []models.Message
	err = facades.Orm().Query().
		Where("chat_id", chat.ID).
		Get(&messages)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, messages)
}

func (r *MessageController) Show(ctx http.Context) http.Response {
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

	var message models.Message
	err = facades.Orm().Query().
		Where("number", ctx.Request().Input("number")).
		Where("chat_id", chat.ID).
		FirstOrFail(&message)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, message)
}

func (r *MessageController) Store(ctx http.Context) http.Response {
	var application models.Application
	err := facades.Orm().Query().Where("token", ctx.Request().Input("token")).FirstOrFail(&application)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}

	var chat models.Chat
	err = facades.Orm().Query().
		Where("number", ctx.Request().Input("number")).
		Where("application_id", application.ID).
		FirstOrFail(&chat)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}

	_, err = facades.Orm().Query().Exec("INSERT INTO messages (number, chat_id, body) SELECT COALESCE(MAX(number), 0)+1, ?, ? FROM messages WHERE chat_id = ?", chat.ID, chat.ID, ctx.Request().Input("body"))
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}
	return ctx.Response().Redirect(http.StatusCreated, "/api/applications/"+ctx.Request().Input("token")+"/chats/"+ctx.Request().Input("number")+"/messages")
}
