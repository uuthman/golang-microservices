package services

import (
	"github.com/uuthman/golang-microservices/src/api/config"
	"github.com/uuthman/golang-microservices/src/api/providers/github_provider"
	"github.com/uuthman/golang-microservices/src/api/domain/github"
	"strings"
	"github.com/uuthman/golang-microservices/src/api/domain/repositories"
	"github.com/uuthman/golang-microservices/src/api/utils/errors"
)

type reposService struct{}

type reposServiceInterface interface{
	CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
}

var (
	RepositoryService reposServiceInterface
)

func init(){
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError){
	input.Name  = strings.TrimSpace(input.Name)
	if input.Name == ""{
		return nil,errors.NewBadRequestError("invalid repository name")
	}

	request := github.CreateRepoRequest{
		Name: input.Name,
		Description:input.Description,
		Private: false,
	}

	response,err := github_provider.CreateRepo(config.GetGithubAccessToken(),request)

	if err != nil {
		
		return nil,errors.NewApiError(err.StatusCode,err.Message)
	}

	result := &repositories.CreateRepoResponse{
		ID: response.ID,
		Name: response.Name,
		Owner: response.Owner.Login,
	}

	return result,nil
}