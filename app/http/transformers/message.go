package transformers

import "chat/app/models"

func MessageResponse(message models.Message) map[string]interface{} {
	return map[string]interface{}{
		"chat_number": message.Chat.Number,
		"number":      message.Number,
		"content":     message.Body,
		"created_at":  message.CreatedAt,
		"updated_at":  message.UpdatedAt,
	}
}

func MessagesCollectionResponse(messages []models.Message) []map[string]interface{} {
	var data []map[string]interface{}
	for _, message := range messages {
		data = append(data, MessageResponse(message))
	}
	return data
}
