package domain

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetUserNoUserFound(t *testing.T){
	user,err := UserDao.GetUser(0)
	assert.Nil(t,user,"we were not expecting a user with id 0")
	assert.NotNil(t,err)
	assert.EqualValues(t,http.StatusNotFound,err.StatusCode)
	assert.EqualValues(t,"not_found",err.Code)
	assert.EqualValues(t,"users 0 was not found",err.Message)
}


func TestGetUserNoError(t *testing.T){
	user,err := UserDao.GetUser(123)
	assert.Nil(t,err)
	assert.NotNil(t,user)
	assert.EqualValues(t,123,user.ID)
	assert.EqualValues(t,"Uthman",user.FirstName)
	assert.EqualValues(t,"Ayinde",user.LastName)
	assert.EqualValues(t,"uthman@email.com",user.Email)

}