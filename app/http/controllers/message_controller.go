package controllers

import (
	"chat/app/http/transformers"
	"chat/app/services"
	workers "chat/pkg/rabbitmq"
	"encoding/json"
	"github.com/streadway/amqp"

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

func (r *MessageController) Search(ctx http.Context) http.Response {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			return
		}
	}(ch)

	msg := workers.Message{
		AppToken: ctx.Request().Input("token"),
		ChatId:   ctx.Request().InputInt("number"),
		Body:     ctx.Request().Input("body"),
	}
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}

	err = ch.Publish(
		"",              // exchange
		"message_queue", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgJSON,
		})
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}

	return ctx.Response().Json(http.StatusCreated, nil)
}
