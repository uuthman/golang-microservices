package github_provider

import (
	"strings"
	"io/ioutil"
	"os"
	"errors"
	"net/http"
	"github.com/uuthman/golang-microservices/src/api/domain/github"
	"github.com/uuthman/golang-microservices/src/api/clients/restclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M){
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestGetAuthorizationHeader(t *testing.T){

	header := getAuthorizationHeader("abc123")

	assert.EqualValues(t,"token abc123",header)
}

func TestCreateRepoErrorRestClient(t *testing.T){
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		URL: "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Err: errors.New("invalid restclient response"),
	})

	response , err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,response)
	assert.NotNil(t,err)
	assert.EqualValues(t,http.StatusInternalServerError,err.StatusCode)
	assert.EqualValues(t,"invalid restclient response",err.Message)


}

func TestCreateRepoInvalidResponseBody(t *testing.T){
	restclient.FlushMockups()
	i,_ := os.Open("fjfj")
	restclient.AddMockup(restclient.Mock{
		URL: "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response:&http.Response{
			StatusCode: http.StatusCreated,
			Body: i,
		},
	})

	response , err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,response)
	assert.NotNil(t,err)
	assert.EqualValues(t,http.StatusInternalServerError,err.StatusCode)
	assert.EqualValues(t,"invalid response json",err.Message)


}

func TestCreateRepoInvalidErrorInterface(t *testing.T){
	restclient.FlushMockups()
	
	restclient.AddMockup(restclient.Mock{
		URL: "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response:&http.Response{
			StatusCode: http.StatusUnauthorized,
			Body: ioutil.NopCloser(strings.NewReader(`{"message":1}`)),
		},
	})

	response , err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,response)
	assert.NotNil(t,err)
	assert.EqualValues(t,http.StatusInternalServerError,err.StatusCode)
	assert.EqualValues(t,"invalid json response body",err.Message)


}

func TestCreateRepoUnautorized(t *testing.T){
	restclient.FlushMockups()
	
	restclient.AddMockup(restclient.Mock{
		URL: "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response:&http.Response{
			StatusCode: http.StatusUnauthorized,
			Body: ioutil.NopCloser(strings.NewReader(`{"message":"requires authentication"}`)),
		},
	})

	response , err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,response)
	assert.NotNil(t,err)
	assert.EqualValues(t,http.StatusUnauthorized,err.StatusCode)
	assert.EqualValues(t,"requires authentication",err.Message)


}

func TestCreateRepoInvalidSuccessReponse(t *testing.T){
	restclient.FlushMockups()
	
	restclient.AddMockup(restclient.Mock{
		URL: "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response:&http.Response{
			StatusCode: http.StatusCreated,
			Body: ioutil.NopCloser(strings.NewReader(`{"id":"123"}`)),
		},
	})

	response , err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,response)
	assert.NotNil(t,err)
	assert.EqualValues(t,http.StatusInternalServerError,err.StatusCode)
	assert.EqualValues(t,"error when trying to unmarshal github create repo response",err.Message)


}

func TestCreateRepoNoError(t *testing.T){
	restclient.FlushMockups()
	
	restclient.AddMockup(restclient.Mock{
		URL: "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response:&http.Response{
			StatusCode: http.StatusCreated,
			Body: ioutil.NopCloser(strings.NewReader(`{"id":123,"name":"tutorial","full_name":"test"}`)),
		},
	})

	response , err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,err)
	assert.NotNil(t,response)
	assert.EqualValues(t,123,response.ID)
	assert.EqualValues(t,"tutorial",response.Name)
	


}