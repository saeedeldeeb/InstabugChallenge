package controllers

import (
	"chat/app/http/transformers"
	"chat/app/models"
	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ApplicationController struct {
	//Dependent services
}

func NewApplicationController() *ApplicationController {
	return &ApplicationController{
		//Inject services
	}
}

func (r *ApplicationController) Index(ctx http.Context) http.Response {
	var applications []models.Application
	err := facades.Orm().Query().Get(&applications)
	if err != nil {
		return nil
	}
	return ctx.Response().Json(http.StatusOK, transformers.ApplicationsCollectionResponse(applications))
}

func (r *ApplicationController) Show(ctx http.Context) http.Response {
	var application models.Application
	err := facades.Orm().Query().Where("token", ctx.Request().Input("token")).FirstOrFail(&application)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.ApplicationResponse(application))
}

func (r *ApplicationController) Store(ctx http.Context) http.Response {
	application := models.Application{
		Name:  ctx.Request().Input("name"),
		Token: uuid.New().String(),
	}
	err := facades.Orm().Query().Create(&application)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}
	return ctx.Response().Json(http.StatusCreated, transformers.ApplicationResponse(application))
}

func (r *ApplicationController) Update(ctx http.Context) http.Response {
	var application models.Application
	err := facades.Orm().Query().Where("token", ctx.Request().Input("token")).FirstOrFail(&application)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, nil)
	}
	application.Name = ctx.Request().Input("name")
	_, err = facades.Orm().Query().Update(&application)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, nil)
	}
	return ctx.Response().Json(http.StatusOK, transformers.ApplicationResponse(application))
}
