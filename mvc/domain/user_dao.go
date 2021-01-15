package domain

import (
	"log"
	"net/http"
	"fmt"
	"github.com/uuthman/golang-microservices/mvc/utils"

)

var (
	users = map[int64]*User{
		123:{ID:123,FirstName:"Uthman",LastName:"Ayinde",Email:"uthman@email.com",},
	}

	UserDao userDaoInterface
)

type userDao struct{}

func init(){
	UserDao = &userDao{}
}


type userDaoInterface interface{
	GetUser(int64) (*User,*utils.ApplicationError)
}

func(u *userDao) GetUser(userId int64) (*User,*utils.ApplicationError){
	log.Println("we're accessing the database")
	
	if user := users[userId]; user != nil{
		return user,nil
	}
	return nil,&utils.ApplicationError{
		Message: fmt.Sprintf("users %v was not found",userId),
		StatusCode:http.StatusNotFound,
		Code: "not_found",
	}
}