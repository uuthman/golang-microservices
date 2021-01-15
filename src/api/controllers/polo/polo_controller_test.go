package polo

import (
	"github.com/uuthman/golang-microservices/src/api/utils/test_utils"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContants(t *testing.T){
	assert.EqualValues(t,"polo",polo)	
}

func TestPolo(t *testing.T){
	response := httptest.NewRecorder()
	request,_ := http.NewRequest(http.MethodGet,"/marco",nil)

	c := test_utils.GetMockedContext(request,response)

	Polo(c)

	assert.EqualValues(t,http.StatusOK,response.Code)
	assert.EqualValues(t,"polo",response.Body.String())

}