package app

import (
	"github.com/uuthman/golang-microservices/src/api/controllers/polo"
	"github.com/uuthman/golang-microservices/src/api/controllers/repositories"
)


func mapUrls(){
	router.GET("/marco",polo.Polo)
	router.POST("/repository",repositories.CreateRepo)
	router.POST("/repositories",repositories.CreateRepos)
}