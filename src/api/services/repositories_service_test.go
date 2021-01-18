package services

import (
	"github.com/uuthman/golang-microservices/src/api/utils/errors"
	"sync"
	"strings"
	"io/ioutil"
	"net/http"
	"github.com/stretchr/testify/assert"
	"github.com/uuthman/golang-microservices/src/api/domain/repositories"
	"os"
	"github.com/uuthman/golang-microservices/src/api/clients/restclient"
	"testing"
)

func TestMain(m *testing.M){
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidInputName(t *testing.T){
	request := repositories.CreateRepoRequest{}

	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t,result)
	assert.NotNil(t,err)
	assert.EqualValues(t,http.StatusBadRequest,err.Status())
	assert.EqualValues(t,"invalid repository name",err.Message())
}

func TestCreateRepoErrorFromGithub(t *testing.T){
	restclient.FlushMockups()

	restclient.AddMockup(
		restclient.Mock{
			URL:"https://api.github.com/user/repos",
			HTTPMethod:http.MethodPost,
			Response: &http.Response{
				StatusCode:http.StatusUnauthorized,
				Body: ioutil.NopCloser(strings.NewReader(`{"message":"requires authentication"}`)),
			},	})

	request := repositories.CreateRepoRequest{Name: "testing"}


	result, err :=	RepositoryService.CreateRepo(request)
	
	assert.Nil(t,result)
	assert.NotNil(t,err)
	assert.EqualValues(t,http.StatusUnauthorized,err.Status())
	assert.EqualValues(t,"requires authentication",err.Message())


}

func TestCreateRepoNoError(t *testing.T){
	restclient.FlushMockups()

	restclient.AddMockup(
		restclient.Mock{
			URL:"https://api.github.com/user/repos",
			HTTPMethod:http.MethodPost,
			Response: &http.Response{
				StatusCode:http.StatusCreated,
				Body: ioutil.NopCloser(strings.NewReader(`{"id":123,"name":"testing","owner":{"login":"uthman"}}`)),
			},	})

	request := repositories.CreateRepoRequest{Name: "testing"}


	result, err :=	RepositoryService.CreateRepo(request)
	
	assert.Nil(t,err)
	assert.NotNil(t,result)
	assert.EqualValues(t,123,result.ID)
	assert.EqualValues(t,"testing",result.Name)
	assert.EqualValues(t,"uthman",result.Owner)

}

func TestCreateRepoConcurrentInvalidRequest(t *testing.T){

	request := repositories.CreateRepoRequest{}
	
	output := make(chan repositories.CreateRepositoriesResult)
	
	service := reposService{}
	go service.createRepoCuncurrent(request,output)

	result :=  <- output

	assert.NotNil(t,result)
	assert.Nil(t,result.Response)
	assert.NotNil(t,result.Error)
	assert.EqualValues(t,http.StatusBadRequest,result.Error.Status())
	assert.EqualValues(t,"invalid repository name",result.Error.Message())



}

func TestCreateRepoConcurrentErrorFromGithub(t *testing.T){
	restclient.FlushMockups()

	restclient.AddMockup(
		restclient.Mock{
			URL:"https://api.github.com/user/repos",
			HTTPMethod:http.MethodPost,
			Response: &http.Response{
				StatusCode:http.StatusUnauthorized,
				Body: ioutil.NopCloser(strings.NewReader(`{"message":"requires authentication"}`)),
			},	})
	
			
	request := repositories.CreateRepoRequest{Name:"testing"}
	
	output := make(chan repositories.CreateRepositoriesResult)
			
	service := reposService{}
	go service.createRepoCuncurrent(request,output)
		
	result :=  <- output
		
	assert.NotNil(t,result)
	assert.Nil(t,result.Response)
	assert.NotNil(t,result.Error)
	assert.EqualValues(t,http.StatusUnauthorized,result.Error.Status())
	assert.EqualValues(t,"requires authentication",result.Error.Message())
				
}

func TestCreateRepoConcurrentNoError(t *testing.T){
	restclient.FlushMockups()

	restclient.AddMockup(
		restclient.Mock{
			URL:"https://api.github.com/user/repos",
			HTTPMethod:http.MethodPost,
			Response: &http.Response{
				StatusCode:http.StatusCreated,
				Body: ioutil.NopCloser(strings.NewReader(`{"id":123,"name":"testing","owner":{"login":"uthman"}}`)),
			},	
		
	})

	request := repositories.CreateRepoRequest{Name:"testing"}
	
	output := make(chan repositories.CreateRepositoriesResult)
			
	service := reposService{}
	go service.createRepoCuncurrent(request,output)
		
	result :=  <- output

	assert.NotNil(t,result.Response)
	assert.Nil(t,result.Error)
	assert.EqualValues(t,123,result.Response.ID)
	assert.EqualValues(t,"testing",result.Response.Name)
	assert.EqualValues(t,"uthman",result.Response.Owner)


}

func TestHandleRepoResult(t *testing.T){
	output := make(chan repositories.CreateReposResponse)
	input := make(chan repositories.CreateRepositoriesResult)

	var wg sync.WaitGroup

	service := reposService{}

	wg.Add(1)
	go func(){
		input <- repositories.CreateRepositoriesResult{
			Error: errors.NewBadRequestError("invalid repo name"),
		}
	}()

	go service.handleRepoResults(&wg,input,output)

	wg.Wait()
	close(input)
	
	result := <- output

	assert.NotNil(t,result)
	assert.EqualValues(t,0,result.StatusCode)
	assert.EqualValues(t,1,len(result.Results))
	assert.EqualValues(t,http.StatusBadRequest,result.Results[0].Error.Status())
	assert.EqualValues(t,"invalid repo name",result.Results[0].Error.Message())
}

func TestCreateReposInvalidRequest(t *testing.T){

	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "  "},
	}

	result,err := RepositoryService.CreateRepos(requests)

	assert.NotNil(t,result)
	assert.Nil(t,err)

	assert.NotNil(t,result)
	assert.EqualValues(t,http.StatusBadRequest,result.StatusCode)
	assert.EqualValues(t,2,len(result.Results))

	assert.Nil(t,result.Results[0].Response)
	assert.EqualValues(t,http.StatusBadRequest,result.Results[0].Error.Status())
	assert.EqualValues(t,"invalid repository name",result.Results[0].Error.Message())

	assert.Nil(t,result.Results[1].Response)
	assert.EqualValues(t,http.StatusBadRequest,result.Results[1].Error.Status())
	assert.EqualValues(t,"invalid repository name",result.Results[1].Error.Message())



}


func TestCreateReposOnSuccessOnFail(t *testing.T){

	restclient.FlushMockups()

	restclient.AddMockup(
		restclient.Mock{
			URL:"https://api.github.com/user/repos",
			HTTPMethod:http.MethodPost,
			Response: &http.Response{
				StatusCode:http.StatusCreated,
				Body: ioutil.NopCloser(strings.NewReader(`{"id":123,"name":"testing","owner":{"login":"uthman"}}`)),
			},	})


	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "testing"},
	}

	result,err := RepositoryService.CreateRepos(requests)

	assert.NotNil(t,result)
	assert.Nil(t,err)

	assert.EqualValues(t,http.StatusPartialContent,result.StatusCode)
	assert.EqualValues(t,2,len(result.Results))

	for _, result := range result.Results {
		
		if result.Error != nil{
			assert.EqualValues(t,http.StatusBadRequest,result.Error.Status())
			assert.EqualValues(t,"invalid repository name",result.Error.Message())
			continue
		}
		assert.EqualValues(t,123,result.Response.ID)
		assert.EqualValues(t,"testing",result.Response.Name)
		assert.EqualValues(t,"uthman",result.Response.Owner)
	}

	



}


func TestCreateReposAllSuccess(t *testing.T){

	restclient.FlushMockups()

	restclient.AddMockup(
		restclient.Mock{
			URL:"https://api.github.com/user/repos",
			HTTPMethod:http.MethodPost,
			Response: &http.Response{
				StatusCode:http.StatusCreated,
				Body: ioutil.NopCloser(strings.NewReader(`{"id":123,"name":"testing","owner":{"login":"uthman"}}`)),
			},	})


	requests := []repositories.CreateRepoRequest{
		{Name: "testing"},
		
	}

	result,err := RepositoryService.CreateRepos(requests)

	assert.NotNil(t,result)
	assert.Nil(t,err)

	assert.EqualValues(t,http.StatusCreated,result.StatusCode)
	assert.EqualValues(t,1,len(result.Results))

	assert.Nil(t,result.Results[0].Error)
	assert.EqualValues(t,123,result.Results[0].Response.ID)
	

	// assert.Nil(t,result.Results[1].Error)
	// assert.EqualValues(t,123,result.Results[1].Response.ID)
	



}
