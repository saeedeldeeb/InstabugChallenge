package transformers

import "chat/app/models"

func ApplicationResponse(app models.Application) map[string]interface{} {
	return map[string]interface{}{
		"name":        app.Name,
		"token":       app.Token,
		"chats_count": app.ChatsCount,
		"created_at":  app.CreatedAt,
		"updated_at":  app.UpdatedAt,
	}
}

func ApplicationsCollectionResponse(apps []models.Application) []map[string]interface{} {
	var response []map[string]interface{}
	for _, app := range apps {
		response = append(response, ApplicationResponse(app))
	}
	return response
}
