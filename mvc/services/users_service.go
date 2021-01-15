package services

import (
	"github.com/uuthman/golang-microservices/mvc/utils"
	"github.com/uuthman/golang-microservices/mvc/domain"
)

type usersServices struct{}

var (
	UsersServices usersServices
)

func (u *usersServices) GetUser(userId int64) (*domain.User,*utils.ApplicationError){
	return domain.UserDao.GetUser(userId)
}