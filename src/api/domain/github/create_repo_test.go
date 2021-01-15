package github

import (
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"testing"
)

func TestCreateRepoRequestAsJson(t *testing.T){
	request := CreateRepoRequest{
		Name: "golang introduction",
		Description: "a golang introduction repository",
		Homepage: "https://github.com",
		Private: true,
		HasIssues:true,
		HasProjects:true,
		HasWiki:true,
	}

	bytes,err := json.Marshal(request)

	assert.Nil(t,err)
	assert.NotNil(t,bytes)
	
	var target CreateRepoRequest

	err = json.Unmarshal(bytes,&target)

	assert.Nil(t,err)
	assert.EqualValues(t,target.Name,request.Name)
	assert.EqualValues(t,target.Private,request.Private)
}