package routes

import (
	"chat/app/http/controllers"
	"os"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
)

func Api() {
	facades.Route().Prefix("api").Group(func(router route.Router) {
		router.Prefix("applications").Group(func(router route.Router) {
			applicationController := controllers.NewApplicationController()

			router.Get("/", applicationController.Index)
			router.Get("/{token}", applicationController.Show)
			router.Post("/", applicationController.Store)
			router.Put("/{token}", applicationController.Update)
		})

		router.Prefix("applications/{token}/chats").Group(func(router route.Router) {
			chatController := controllers.NewChatController()

			router.Get("/", chatController.Index)
			router.Get("/{number}", chatController.Show)
			router.Post("/", chatController.Store)
		})

		router.Prefix("applications/{token}/chats/{number}/messages").Group(func(router route.Router) {
			messageController := controllers.NewMessageController()

			router.Get("/search", messageController.Search)
			router.Get("/", messageController.Index)
			router.Get("/{msg_number}", messageController.Show)
			router.Post("/", messageController.Store)
		})
	})

	facades.Route().Get("/ping", func(ctx http.Context) http.Response {
		hostname, _ := os.Hostname() // Get container hostname
		return ctx.Response().Json(200, map[string]interface{}{
			"message":  "pong",
			"ip":       ctx.Request().Ip(),
			"instance": hostname,
		})
	})
}
