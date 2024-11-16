package listeners

import (
	"chat/app/models"
	elasticsearch "chat/pkg"
	"context"
	"encoding/json"
	"github.com/goravel/framework/facades"
	"strconv"

	"github.com/goravel/framework/contracts/event"
)

type IndexMessageInElasticSearch struct {
}

func (receiver *IndexMessageInElasticSearch) Signature() string {
	return "index_message_in_elastic_search"
}

func (receiver *IndexMessageInElasticSearch) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     false,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *IndexMessageInElasticSearch) Handle(args ...any) error {
	messageJSON := args[0].(string)
	var message models.Message
	err := json.Unmarshal([]byte(messageJSON), &message)
	if err != nil {
		return err
	}

	// TODO: get configs from env
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "elasticsearch",
		Password:  "password123",
	})
	if err != nil {
		return err
	}

	err = es.Index(context.Background(), "messages", strconv.Itoa(int(message.ID)), message)
	if err != nil {
		facades.Log().Error(err)
		return err
	}

	return nil
}
