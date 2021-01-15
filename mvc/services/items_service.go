package services


import (
	"net/http"
	"github.com/uuthman/golang-microservices/mvc/utils"
	"github.com/uuthman/golang-microservices/mvc/domain"
)

type itemsService struct{}

var (
	ItemsService itemsService
)



func GetItem(itemId string)(*domain.Item,*utils.ApplicationError){
	return nil,&utils.ApplicationError{
		Message: "implement me",
		StatusCode: http.StatusInternalServerError,
	}
}