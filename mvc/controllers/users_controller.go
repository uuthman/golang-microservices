package controllers

import (
	"net/http"
	"github.com/uuthman/golang-microservices/mvc/utils"
	"github.com/uuthman/golang-microservices/mvc/services"
	"strconv"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context){
	
	userId,err := strconv.ParseInt(c.Param("user_id"),10,64)
	if err != nil{

		apiErr := &utils.ApplicationError{
			Message: "user_id must be a number",
			StatusCode:http.StatusBadGateway,
			Code:"bad_request",
		}

		utils.RespondError(c,apiErr)
		return
	}

	user,apiErr := services.UsersServices.GetUser(userId)
	if apiErr != nil{	
		utils.RespondError(c,apiErr)
		return
	}

	utils.Respond(c,http.StatusOK,user)

}