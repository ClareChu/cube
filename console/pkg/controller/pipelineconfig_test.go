package controller

import (
	"fmt"
	"hidevops.io/cube/console/pkg/aggregate/mocks"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"hidevops.io/hiboot/pkg/utils/io"
	"net/http"
	"testing"
	"time"
)

func genJwtToken(timeout int64) (token string) {
	jwtToken := jwt.NewJwtToken(&jwt.Properties{
		PrivateKeyPath: "config/ssl/app.rsa",
		PublicKeyPath:  "config/ssl/app.rsa.pub",
	})
	pt, err := jwtToken.Generate(jwt.Map{
		"username": "johndoe",
		"password": "PA$$W0RD",
	}, timeout, time.Second)
	if err == nil {
		token = fmt.Sprintf("Bearer %v", pt)
	}
	return
}

func TestPipelineConfig(t *testing.T) {

	pipelineConfigAggregate := new(mocks.PipelineConfigAggregate)
	secrete := new(mocks.SecretAggregate)
	pipelineConfig := newPipelineConfigController(pipelineConfigAggregate, secrete)
	s := &command.Secret{
		Username: "johndoe",
		Password: "PA$$W0RD",
	}
	secrete.On("Create", s).Return(nil)
	cmd := &command.PipelineStart{}
	pipelineConfigAggregate.On("StartPipelineConfig", cmd).Return(nil, nil)
	pipelineConfigAggregate.On("Get", "a", "b").Return(nil, nil)
	app := web.NewTestApp(t, pipelineConfig).
		SetProperty("kube.serviceHost", "test").
		Run(t)
	log.SetLevel(log.DebugLevel)
	log.Println(io.GetWorkDir())
	token := genJwtToken(100)
	time.Sleep(2 * time.Second)
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Post("/pipelineConfig").WithHeader("Authorization", token).WithJSON(&command.PipelineConfigTemplate{}).Expect().Status(http.StatusOK)
	})
	time.Sleep(2 * time.Second)
	t.Run("should pass with jwt token", func(t *testing.T) {
		app.Get("/pipelineConfig/name/a/namespace/b").WithHeader("Authorization", token).Expect().Status(http.StatusOK)
	})
}
