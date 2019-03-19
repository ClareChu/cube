package controller

import (
	"errors"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"hidevops.io/mio/pkg/auth/service"
	"net/http"
	"net/url"
	"os"
	"time"
)

type OAuth struct {
}

type oauthController struct {
	at.RestController
	token            jwt.Token
	oauthService     service.OauthService
	loginService     service.LoginService
	sessionInterface service.SessionInterface
}

func init() {
	app.Register(newOauthController)
}

func newOauthController(token jwt.Token, oauthService service.OauthService, loginService service.LoginService, sessionInterface service.SessionInterface) *oauthController {
	return &oauthController{
		token:          token,
		oauthService:   oauthService,
		loginService:   loginService,
		sessionInterface: sessionInterface,
	}
}

func (o *oauthController) GetUrl() (response model.Response, err error) {
	log.Debug("gitlab get oauth2 url")
	response = new(model.BaseResponse)
	a := &service.Auth{
		AuthURL:       os.Getenv("SCM_URL"),
		ApplicationId: service.ApplicationId,
		CallbackUrl:   service.CallbackUrl,
	}
	url := o.oauthService.GetAuthURL(a)
	responseBody := map[string]interface{}{
		"authUrl": &url,
	}
	response.SetData(responseBody)
	return
}

func (o *oauthController) GetByCode(code string) (response model.Response, err error) {
	log.Debug("gitlab oauth2 login ")
	//TODO 通过code  获取用户的 access token  和失效时间  通过时间来获取用户信息
	response = new(model.BaseResponse)
	s := url.QueryEscape(service.CallbackUrl)
	session := service.NewClient(service.BaseUrl, service.AccessTokenUrl, service.ApplicationId, s, service.Secret)
	resp, err := o.sessionInterface.GetAccessToken(session, code)
	if err != nil || resp.AccessToken == "" {
		response.SetCode(http.StatusUnauthorized)
		log.Error("login session get accessToken error", err)
		return response, errors.New("login session get accessToken error")
	}

	//TODO 通过accessToken 获取用户信息
	accessToken := resp.AccessToken
	log.Debugf("accessToken : %v", accessToken)
	u, err := o.loginService.GetUser(accessToken)
	token, err := o.token.Generate(jwt.Map{
		"username":     u.Name,
		"access_token": accessToken,
		"uid":          u.ID,
	}, 1000, time.Hour)
	if err != nil {
		response.SetCode(http.StatusNonAuthoritativeInfo)
		return
	}

	data := map[string]interface{}{
		"token": token,
		"user":  u,
	}
	response.SetData(data)
	return
}
