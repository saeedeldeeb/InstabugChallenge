package routes

import (
	"chat/app/http/controllers"
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
	})
}
