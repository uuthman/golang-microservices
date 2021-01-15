package restclient

import (
	"fmt"
	"errors"
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	enabledMocks = false
	mocks = make(map[string]*Mock)
)

type Mock struct{
	URL string
	HTTPMethod string
	Response *http.Response
	Err error
}

func getMockId(httpMethod string,url string) string {
	return fmt.Sprintf("%s_%s",httpMethod,url)
}

func StartMockups(){
	enabledMocks = true
}

func StopMockups(){
	enabledMocks = false
}

func AddMockup(mock Mock){
	mocks[getMockId(mock.HTTPMethod,mock.URL)] = &mock
}

func FlushMockups(){
	mocks = make(map[string]*Mock)
	
}

func Post(url string,body interface{},headers http.Header) (*http.Response,error){
	if enabledMocks{
		mock := mocks[getMockId(http.MethodPost,url)]
		if mock == nil{
			return nil , errors.New("no mockup found for given request")
		}
		return mock.Response,mock.Err
	}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req,err := http.NewRequest(http.MethodPost,url,bytes.NewReader(jsonBytes))
	req.Header = headers

	client := http.Client{}
	return client.Do(req)
}