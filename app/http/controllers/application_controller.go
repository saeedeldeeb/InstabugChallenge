package controllers

import (
	"chat/app/http/transformers"
	"chat/app/services"
	"github.com/goravel/framework/contracts/http"
)

type ApplicationController struct {
	//Dependent services
	applicationService services.Application
}

func NewApplicationController() *ApplicationController {
	return &ApplicationController{
		//Inject services
		applicationService: services.NewApplicationService(),
	}
}

func (r *ApplicationController) Index(ctx http.Context) http.Response {
	applications, err := r.applicationService.GetApplications()
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.ApplicationsCollectionResponse(applications))
}

func (r *ApplicationController) Show(ctx http.Context) http.Response {
	application, err := r.applicationService.GetApplicationByToken(ctx.Request().Input("token"))
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.ApplicationResponse(application))
}

func (r *ApplicationController) Store(ctx http.Context) http.Response {
	application, err := r.applicationService.CreateApplication(ctx.Request().Input("name"))
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}
	return ctx.Response().Json(http.StatusCreated, transformers.ApplicationResponse(application))
}

func (r *ApplicationController) Update(ctx http.Context) http.Response {
	application, err := r.applicationService.UpdateApplication(ctx.Request().Input("token"), ctx.Request().Input("name"))
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.ApplicationResponse(application))
}
