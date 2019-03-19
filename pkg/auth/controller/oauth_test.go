package controller

import (
	"errors"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"hidevops.io/hiboot/pkg/utils/io"
	"hidevops.io/hioak/starter/scm"
	"hidevops.io/mio/pkg/auth/service"
	"hidevops.io/mio/pkg/auth/service/mocks"
	"net/http"
	"net/url"
	"testing"
)

func TestOauthControllerGetUrl(t *testing.T) {
	io.EnsureWorkDir(1, "config/application.yml")
	jwtToken := jwt.NewJwtToken(&jwt.Properties{
		PrivateKeyPath: "config/ssl/app.rsa",
		PublicKeyPath:  "config/ssl/app.rsa.pub",
	})
	authService := new(mocks.OauthService)
	loginService := new(mocks.LoginService)
	sessionService := new(mocks.SessionInterface)
	controller := newOauthController(jwtToken, authService, loginService, sessionService)
	application := web.NewTestApp(t, controller).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	authService.On("GetAuthURL").Return(nil, "")
	t.Run("should pass with jwt token", func(t *testing.T) {
		application.Get("/oauth/url").Expect().Status(http.StatusOK)
	})

	authService.On("GetAuthURL").Return(errors.New("not found"), "")
	t.Run("should pass with jwt token", func(t *testing.T) {
		application.Get("/oauth/url").Expect().Status(http.StatusOK)
	})
	s := url.QueryEscape(service.CallbackUrl)
	session := service.NewClient(service.OauthUrl, service.AccessTokenUrl, service.ApplicationId, s, service.Secret)
	re := &service.SessionResponse{
		AccessToken: "q",
	}
	sessionService.On("GetAccessToken", session, "aa").Return(re, nil)

	loginService.On("GetUser", "q").Return(&scm.User{}, nil)
	t.Run("should pass with jwt token", func(t *testing.T) {
		application.Get("/oauth/code/aa").Expect().Status(http.StatusOK)
	})
}
