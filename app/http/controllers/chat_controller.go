package controllers

import (
	"chat/app/http/transformers"
	"chat/app/services"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ChatController struct {
	//Dependent services
	chatService services.Chat
}

func NewChatController() *ChatController {
	return &ChatController{
		//Inject services
		chatService: services.NewChatService(),
	}
}

func (r *ChatController) Index(ctx http.Context) http.Response {
	chats, err := r.chatService.GetChats(ctx.Request().Input("token"))
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.ChatsCollectionResponse(chats))
}

func (r *ChatController) Show(ctx http.Context) http.Response {
	chat, err := r.chatService.GetChatByNumber(ctx.Request().Input("token"), ctx.Request().InputInt("number"))
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.ChatResponse(chat))
}

func (r *ChatController) Store(ctx http.Context) http.Response {
	chat, err := r.chatService.CreateChat(ctx.Request().Input("token"))
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}
	return ctx.Response().Json(http.StatusCreated, transformers.ChatResponse(chat))
}
