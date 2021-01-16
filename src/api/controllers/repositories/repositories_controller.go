package repositories

import (
	"net/http"
	"github.com/uuthman/golang-microservices/src/api/services"
	"github.com/uuthman/golang-microservices/src/api/utils/errors"
	"github.com/uuthman/golang-microservices/src/api/domain/repositories"
	"github.com/gin-gonic/gin"
)

func CreateRepo (c *gin.Context){
	var request repositories.CreateRepoRequest

	if err := c.ShouldBindJSON(&request); err != nil{
		apiErr := errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.Status(),apiErr)
		return 
	}

	result, err := services.RepositoryService.CreateRepo(request)
	if err != nil {
		c.JSON(err.Status(),err)
		return 
	}

	c.JSON(http.StatusCreated,result)

}

func CreateRepos (c *gin.Context){
	var request []repositories.CreateRepoRequest

	if err := c.ShouldBindJSON(&request); err != nil{
		apiErr := errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.Status(),apiErr)
		return 
	}

	result, err := services.RepositoryService.CreateRepos(request)
	if err != nil {
		c.JSON(err.Status(),err)
		return 
	}

	c.JSON(result.StatusCode,result)

}
