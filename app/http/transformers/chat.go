package transformers

import "chat/app/models"

func ChatResponse(chat models.Chat) map[string]interface{} {
	return map[string]interface{}{
		"number":         chat.Number,
		"messages_count": chat.MessagesCount,
		"application_id": chat.ApplicationId,
		"created_at":     chat.CreatedAt,
		"updated_at":     chat.UpdatedAt,
	}
}

func ChatsCollectionResponse(chats []models.Chat) []map[string]interface{} {
	var response []map[string]interface{}
	for _, chat := range chats {
		response = append(response, ChatResponse(chat))
	}
	return response
}
