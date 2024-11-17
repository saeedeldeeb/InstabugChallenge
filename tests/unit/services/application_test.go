package tests

import (
	"chat/app/models"
	"chat/app/services"
	"errors"
	"testing"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/testing/mock"
	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
)

func TestApplicationService_Create(t *testing.T) {
	mockFactory := mock.Factory()
	mockOrm := mockFactory.Orm()
	mockQuery := mockFactory.OrmQuery()

	// Setup mock chain
	mockOrm.On("Query").Return(mockQuery)

	// Mock Create operation
	mockQuery.On("Create", testify.Anything).Return(func(app interface{}) error {
		// Simulate ID assignment like a real DB would do
		app.(*models.Application).ID = 1
		return nil
	}).Once()

	// Create service instance
	service := services.NewApplicationService()
	app, err := service.CreateApplication("Test App")

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, uint(1), app.ID)
	assert.Equal(t, "Test App", app.Name)

	// Verify mocks
	mockOrm.AssertExpectations(t)
	mockQuery.AssertExpectations(t)
}

func TestApplicationService_Find(t *testing.T) {
	mockFactory := mock.Factory()
	mockOrm := mockFactory.Orm()
	mockQuery := mockFactory.OrmQuery()

	// Setup mock chain
	mockOrm.On("Query").Return(mockQuery)

	// Mock Where and First operations
	mockQuery.On("Where", "token", "test-token").Return(mockQuery).Once()
	mockQuery.On("FirstOrFail", testify.Anything).Return(func(app interface{}) error {
		// Simulate database result
		result := app.(*models.Application)
		result.ID = 1
		result.Name = "Test App"
		return nil
	}).Once()

	// Create service instance
	service := services.NewApplicationService()
	app, err := service.GetApplicationByToken("test-token")

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, uint(1), app.ID)
	assert.Equal(t, "Test App", app.Name)

	// Verify mocks
	mockOrm.AssertExpectations(t)
	mockQuery.AssertExpectations(t)
}

func TestApplicationService_List(t *testing.T) {
	mockFactory := mock.Factory()
	mockOrm := mockFactory.Orm()
	mockQuery := mockFactory.OrmQuery()

	// Setup mock chain
	mockOrm.On("Query").Return(mockQuery)

	// Mock Find operation
	mockQuery.On("Get", testify.Anything).Return(func(apps interface{}) error {
		// Simulate database results
		result := apps.(*[]models.Application)
		*result = []models.Application{
			{Model: orm.Model{ID: 1}, Name: "App 1"},
			{Model: orm.Model{ID: 2}, Name: "App 2"},
		}
		return nil
	}).Once()

	// Create service instance
	service := services.NewApplicationService()
	apps, err := service.GetApplications()

	// Assertions
	assert.Nil(t, err)
	assert.Len(t, apps, 2)
	assert.Equal(t, uint(1), apps[0].ID)
	assert.Equal(t, "App 1", apps[0].Name)
	assert.Equal(t, uint(2), apps[1].ID)
	assert.Equal(t, "App 2", apps[1].Name)

	// Verify mocks
	mockOrm.AssertExpectations(t)
	mockQuery.AssertExpectations(t)
}

// Test error cases
func TestApplicationService_Create_Error(t *testing.T) {
	mockFactory := mock.Factory()
	mockOrm := mockFactory.Orm()
	mockQuery := mockFactory.OrmQuery()

	// Setup mock chain
	mockOrm.On("Query").Return(mockQuery)

	// Mock Create operation with error
	mockQuery.On("Create", testify.Anything).Return(errors.New("database error")).Once()

	// Create service instance
	service := services.NewApplicationService()
	app, err := service.CreateApplication("Test App")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, models.Application{}, app)
	assert.Equal(t, "database error", err.Error())

	// Verify mocks
	mockOrm.AssertExpectations(t)
	mockQuery.AssertExpectations(t)
}
