package controllers

import (
	"chat/app/http/transformers"
	"chat/app/services"
	"github.com/goravel/framework/contracts/http"
)

type MessageController struct {
	//Dependent services
	msgService services.Message
}

func NewMessageController() *MessageController {
	return &MessageController{
		//Inject services
		msgService: services.NewMessageService(),
	}
}

func (r *MessageController) Index(ctx http.Context) http.Response {
	message, err := r.msgService.GetMessages(ctx.Request().Input("token"), ctx.Request().Input("number"))
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.MessagesCollectionResponse(message))
}

func (r *MessageController) Show(ctx http.Context) http.Response {
	message, err := r.msgService.GetMessageByNumber(ctx.Request().Input("token"), ctx.Request().Input("number"), ctx.Request().InputInt("msg_number"))
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.MessageResponse(message))
}

func (r *MessageController) Store(ctx http.Context) http.Response {
	message, err := r.msgService.CreateMessage(ctx.Request().Input("token"), ctx.Request().Input("number"), ctx.Request().Input("body"))
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}
	return ctx.Response().Json(http.StatusCreated, transformers.MessageResponse(message))
}
