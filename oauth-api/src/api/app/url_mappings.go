package app

import (
	"github.com/uuthman/golang-microservices/oauth-api/src/api/controllers/oauth"
	"github.com/uuthman/golang-microservices/src/api/controllers/polo"
)


func mapUrls(){
	router.GET("/marco",polo.Polo)
	router.POST("/ouath/access_token",oauth.CreateAccessToken)
	router.GET("/oauth/access_token/:token_id",oauth.GetAccessToken)
}