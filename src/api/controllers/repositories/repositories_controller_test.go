package repositories

import (
	"github.com/uuthman/golang-microservices/src/api/utils/test_utils"
	"encoding/json"
	"github.com/uuthman/golang-microservices/src/api/domain/repositories"
	"io/ioutil"
	"os"
	"github.com/uuthman/golang-microservices/src/api/clients/restclient"
	"github.com/uuthman/golang-microservices/src/api/utils/errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"net/http"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M){
	restclient.StartMockups()
	os.Exit(m.Run())
}
func TestCreateRepoInvalidJsonRequest(t *testing.T){
	response := httptest.NewRecorder()
	c,_ := gin.CreateTestContext(response)

	request,_ := http.NewRequest(http.MethodPost,"/repositories",strings.NewReader(``))

	c.Request = request

	CreateRepo(c)

	assert.EqualValues(t,http.StatusBadRequest,response.Code)

	apiErr,err := errors.NewApiErrFromBytes(response.Body.Bytes())

	assert.Nil(t,err)
	assert.NotNil(t,apiErr)
	assert.EqualValues(t,http.StatusBadRequest,apiErr.Status())
	assert.EqualValues(t,"invalid json body",apiErr.Message())


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

			
	response := httptest.NewRecorder()
	request,_ := http.NewRequest(http.MethodPost,"/repositories",strings.NewReader(`{"name":"testing"}`))

	c := test_utils.GetMockedContext(request,response)

	CreateRepo(c)

	assert.EqualValues(t,http.StatusUnauthorized,response.Code)

	apiErr,err := errors.NewApiErrFromBytes(response.Body.Bytes())

	assert.Nil(t,err)
	assert.NotNil(t,apiErr)
	assert.EqualValues(t,http.StatusUnauthorized,apiErr.Status())
	assert.EqualValues(t,"requires authentication",apiErr.Message())


}

func TestCreateRepoNoError(t *testing.T){

	restclient.FlushMockups()

	restclient.AddMockup(
		restclient.Mock{
			URL:"https://api.github.com/user/repos",
			HTTPMethod:http.MethodPost,
			Response: &http.Response{
				StatusCode:http.StatusCreated,
				Body: ioutil.NopCloser(strings.NewReader(`{"id":123}`)),
			},	})

			
	response := httptest.NewRecorder()

	request,_ := http.NewRequest(http.MethodPost,"/repositories",strings.NewReader(`{"name":"testing"}`))
	c := test_utils.GetMockedContext(request,response)
	CreateRepo(c)

	assert.EqualValues(t,http.StatusCreated,response.Code)


	var result repositories.CreateRepoResponse

	err := json.Unmarshal(response.Body.Bytes(),&result)

	assert.Nil(t,err)
	assert.EqualValues(t,123,result.ID)
}