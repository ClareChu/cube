package controller

import (
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"hidevops.io/hiboot/pkg/utils/io"
	"hidevops.io/hioak/starter/scm"
	"hidevops.io/mio/pkg/auth/service"
	"hidevops.io/mio/pkg/auth/service/mocks"
	"net/http"
	"net/url"
	"os"
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
	a := &service.Auth{
		AuthURL:       os.Getenv("SCM_URL"),
		ApplicationId: service.ApplicationId,
		CallbackUrl:   service.CallbackUrl,
	}
	authService.On("GetAuthURL", a).Return("")
	t.Run("should pass with jwt token", func(t *testing.T) {
		application.Get("/oauth/url").Expect().Status(http.StatusOK)
	})

	os.Setenv("SCM_URL", "https://github.com")
	s := url.QueryEscape(service.CallbackUrl)
	session := service.NewClient(service.BaseUrl, service.AccessTokenUrl, service.ApplicationId, s, service.Secret)
	re := &service.SessionResponse{
		AccessToken: "q",
	}
	sessionService.On("GetAccessToken", session, "aa").Return(re, nil)

	loginService.On("GetUser", "https://github.com", "q").Return(&scm.User{}, nil)
	t.Run("should pass with jwt token", func(t *testing.T) {
		application.Get("/oauth/code/aa").Expect().Status(http.StatusOK)
	})
}
