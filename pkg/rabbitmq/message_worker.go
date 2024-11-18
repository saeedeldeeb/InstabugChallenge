package workers

import (
	"chat/app/services"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

type MessageWorker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// Message represents the structure of your message
type Message struct {
	AppToken string `json:"app_token"`
	ChatId   int    `json:"chat_id"`
	Body     string `json:"body"`
}

func NewMessageWorker(amqpURL, queueName string) (*MessageWorker, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare a queue
	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	return &MessageWorker{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

func (w *MessageWorker) Start() error {
	msgs, err := w.channel.Consume(
		w.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Process the message
			err := w.processMessage(d)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				// Nack the message if processing failed
				err := d.Nack(false, true)
				if err != nil {
					return
				}
				continue
			}
			// Acknowledge the message
			err = d.Ack(false)
			if err != nil {
				return
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

func (w *MessageWorker) processMessage(d amqp.Delivery) error {
	var msg Message
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message: %v", err)
	}

	// Handle the message here
	log.Printf("Received message: %+v", msg)

	// Create a new message using message service
	messageService := services.NewMessageService()
	_, err = messageService.CreateMessage(msg.AppToken, strconv.Itoa(msg.ChatId), msg.Body)

	return nil
}

func (w *MessageWorker) Close() {
	if w.channel != nil {
		err := w.channel.Close()
		if err != nil {
			return
		}
	}
	if w.conn != nil {
		err := w.conn.Close()
		if err != nil {
			return
		}
	}
}
