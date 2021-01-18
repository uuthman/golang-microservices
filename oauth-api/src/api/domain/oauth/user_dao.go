package oauth

import (
	"github.com/uuthman/golang-microservices/src/api/utils/errors"
)

const (
	queryUserByUsernameAndPassword = ""
)

var(
	users = map[string]*User{
		"uth": {
			ID:123,
			Username:"uthman",
		},	
	}
)

func GetUserByUsernameAndPassword(username string,password string) (*User,errors.ApiError){

	user := users[username]
	if user == nil{
		return nil,errors.NewNotFoundError("no user found with given parameter")
	}
	return user, nil
}