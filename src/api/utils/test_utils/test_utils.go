package test_utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func GetMockedContext(request *http.Request, response *httptest.ResponseRecorder) *gin.Context{
	c,_ := gin.CreateTestContext(response)
	c.Request = request
	return c
}