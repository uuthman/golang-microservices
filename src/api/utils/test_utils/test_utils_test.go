package test_utils

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/http"
	"testing"
)

func TestGetMockedContext(t *testing.T){
	request,err := http.NewRequest(http.MethodGet,"http://localhost:123/some",nil)
	assert.Nil(t,err)
	response := httptest.NewRecorder()
	request.Header = http.Header{"X-Mock":{"true"}}

	c := GetMockedContext(request,response)

	assert.EqualValues(t,http.MethodGet,c.Request.Method)
	assert.EqualValues(t,"123",c.Request.URL.Port())
	assert.EqualValues(t,"/some",c.Request.URL.Path)
	assert.EqualValues(t,"http",c.Request.URL.Scheme)
	assert.EqualValues(t,1,len(c.Request.Header))
	assert.EqualValues(t,"true",c.GetHeader("x-mock"))
	assert.EqualValues(t,"true",c.GetHeader("X-Mock"))


}