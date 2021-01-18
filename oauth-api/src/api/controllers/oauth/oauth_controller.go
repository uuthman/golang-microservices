package oauth

import (
	"net/http"
	"github.com/uuthman/golang-microservices/oauth-api/src/api/services"
	"github.com/uuthman/golang-microservices/src/api/utils/errors"
	"github.com/uuthman/golang-microservices/oauth-api/src/api/domain/oauth"
	"github.com/gin-gonic/gin"
)

func CreateAccessToken(c *gin.Context){
	var request oauth.AccessTokenRequest

	if err := c.ShouldBindJSON(request); err != nil{
		apiErr := errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.Status(),apiErr)
	}

	token,err := services.OathService.CreateAccessToken(request)
	if err != nil{
		c.JSON(err.Status(),err)
		return
	}

	c.JSON(http.StatusCreated,token)
}

func GetAccessToken(c *gin.Context){
	tokenID := c.Param("token_id")
	token,err := services.OathService.GetAccessToken(tokenID)

	if err != nil{
		c.JSON(err.Status(),err)
	}

	c.JSON(http.StatusOK,token)
}