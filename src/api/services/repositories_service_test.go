package services

import (
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
