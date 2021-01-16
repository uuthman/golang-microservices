package services

import (
	"net/http"
	"sync"
	"github.com/uuthman/golang-microservices/src/api/config"
	"github.com/uuthman/golang-microservices/src/api/providers/github_provider"
	"github.com/uuthman/golang-microservices/src/api/domain/github"

	"github.com/uuthman/golang-microservices/src/api/domain/repositories"
	"github.com/uuthman/golang-microservices/src/api/utils/errors"
)

type reposService struct{}

type reposServiceInterface interface{
	CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

var (
	RepositoryService reposServiceInterface
)

func init(){
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError){
	
	if err := input.Validate(); err != nil{
		return nil,err
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

func (s *reposService) CreateRepos(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError){

	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	var wg sync.WaitGroup
	go s.handleRepoResults(&wg,input,output)

	for _, current := range requests{
		wg.Add(1)
		go s.createRepoCuncurrent(current,input)
	}

	wg.Wait()
	close(input)

	result := <- output

	successCreation := 0

	for _,current := range result.Results{
		if current.Response != nil{
			successCreation++
		}
	}

	if (successCreation == 0){
		result.StatusCode = result.Results[0].Error.Status()
	}else if successCreation == len(requests){
		result.StatusCode = http.StatusCreated
	}else{
		result.StatusCode = http.StatusPartialContent
	}

	return result,nil
}

func (s *reposService) handleRepoResults(wg *sync.WaitGroup,input chan repositories.CreateRepositoriesResult,output chan repositories.CreateReposResponse){

	var results repositories.CreateReposResponse

	for incomingResult := range input{
		repoResult := repositories.CreateRepositoriesResult{
			Response: incomingResult.Response,
			Error: incomingResult.Error,
		}
		results.Results = append(results.Results,repoResult)
		wg.Done()
	}

	output <- results

}

func (s *reposService) createRepoCuncurrent(input repositories.CreateRepoRequest, output chan repositories.CreateRepositoriesResult){

	if err := input.Validate(); err != nil{
		output <- repositories.CreateRepositoriesResult{Error:err}
		return
	}

	result,err := s.CreateRepo(input)

	if err != nil{
		output <- repositories.CreateRepositoriesResult{Error: err}
		return 
	}
	output <- repositories.CreateRepositoriesResult{Response: result}
}