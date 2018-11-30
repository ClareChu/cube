package controller

import (
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"hidevops.io/hiboot/pkg/utils/io"
	"hidevops.io/mio/pkg/auth/service"
	"hidevops.io/mio/pkg/auth/service/mocks"
	"net/http"
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
	controller := newOauthController(jwtToken, authService, loginService)
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
}
