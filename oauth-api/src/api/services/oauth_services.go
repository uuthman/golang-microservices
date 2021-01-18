package services

import (
	"time"
	"github.com/uuthman/golang-microservices/src/api/utils/errors"
	"github.com/uuthman/golang-microservices/oauth-api/src/api/domain/oauth"
)

type oauthService struct{}

type oauthServiceInterface interface{
	CreateAccessToken(request oauth.AccessTokenRequest) (*oauth.AccessToken,errors.ApiError)
	GetAccessToken(accessToken string) (*oauth.AccessToken,errors.ApiError)
}

var (
	OathService oauthServiceInterface
)

func init(){
	OathService = &oauthService{}
}

func (s *oauthService) CreateAccessToken(request oauth.AccessTokenRequest) (*oauth.AccessToken,errors.ApiError){
	if err := request.Validate(); err != nil{
		return nil,err
	}

	user, err := oauth.GetUserByUsernameAndPassword(request.Username,request.Password)
	if err != nil{
		return nil,err
	}

	token := oauth.AccessToken{
		UserID: user.ID,
		Expires: time.Now().UTC().Add(24 * time.Hour).Unix(),
	}

	if err := token.Save(); err != nil{
		return nil,err
	}

	return &token,nil
}

func (s *oauthService) GetAccessToken(accessToken string) (*oauth.AccessToken,errors.ApiError){
	token,err := oauth.GetAccessTokenByToken(accessToken)
	if err != nil{
		return nil,err
	}

	if token.IsExpired(){
		return nil, errors.NewNotFoundError("no access token found with given parameter")
	}
	return token,nil
}
