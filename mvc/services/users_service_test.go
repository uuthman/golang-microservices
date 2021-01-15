package services

import (
	"github.com/uuthman/golang-microservices/mvc/domain"
	"github.com/uuthman/golang-microservices/mvc/utils"
	"net/http"
	"github.com/stretchr/testify/assert"
	
	"testing"
)

var (
	userDaoMock usersDaoMock
	getUserFunction func(userId int64) (*domain.User,*utils.ApplicationError)
)

type usersDaoMock struct{}

func init(){
	domain.UserDao = &usersDaoMock{}
}

//IMPLEMENTS THE userDaoInterface
func (m *usersDaoMock) GetUser(userId int64) (*domain.User,*utils.ApplicationError){
	return getUserFunction(userId)
}


func TestGetUserNotFoundInDatabase(t *testing.T){

	getUserFunction =  func (userId int64) (*domain.User,*utils.ApplicationError){
		return nil,&utils.ApplicationError{
			Message:"users 0 was not found",
			StatusCode: http.StatusNotFound,
		}
	}
	user,err := UsersServices.GetUser(0)

	assert.Nil(t,user)
	assert.NotNil(t,err)
	assert.EqualValues(t,http.StatusNotFound,err.StatusCode)
	assert.EqualValues(t,"users 0 was not found",err.Message)

}

func TestGetUserNoError(t *testing.T){
	getUserFunction =  func (userId int64) (*domain.User,*utils.ApplicationError){
		return &domain.User{
			ID:123,
		},nil
	}

	user,err := UsersServices.GetUser(123)

	assert.Nil(t,err)
	assert.NotNil(t,user)
	assert.EqualValues(t,123,user.ID)

}