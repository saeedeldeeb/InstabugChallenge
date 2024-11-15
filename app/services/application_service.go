package services

import (
	"chat/app/models"
	"github.com/google/uuid"
	"github.com/goravel/framework/facades"
)

type Application interface {
	GetApplications() ([]models.Application, error)
	GetApplicationByToken(token string) (models.Application, error)
	CreateApplication(name string) (models.Application, error)
	UpdateApplication(token, name string) (models.Application, error)
}

type ApplicationService struct {
	//Dependent services
}

func NewApplicationService() *ApplicationService {
	return &ApplicationService{
		//Inject services
	}
}

func (r *ApplicationService) GetApplications() ([]models.Application, error) {
	var applications []models.Application
	err := facades.Orm().Query().Get(&applications)
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *ApplicationService) GetApplicationByToken(token string) (models.Application, error) {
	var application models.Application
	err := facades.Orm().Query().Where("token", token).FirstOrFail(&application)
	if err != nil {
		return models.Application{}, err
	}
	return application, nil
}

func (r *ApplicationService) CreateApplication(name string) (models.Application, error) {
	application := models.Application{
		Name:  name,
		Token: uuid.New().String(),
	}
	err := facades.Orm().Query().Create(&application)
	if err != nil {
		return models.Application{}, err
	}
	return application, nil
}

func (r *ApplicationService) UpdateApplication(token, name string) (models.Application, error) {
	var application models.Application
	err := facades.Orm().Query().Where("token", token).FirstOrFail(&application)
	if err != nil {
		return models.Application{}, err
	}
	application.Name = name
	_, err = facades.Orm().Query().Update(&application)
	if err != nil {
		return models.Application{}, err
	}
	return application, nil
}
