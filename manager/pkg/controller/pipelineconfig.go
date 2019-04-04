package controller

import (
	"github.com/prometheus/common/log"
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/starter/jwt"
	"strings"
)

type PipelineConfigController struct {
	at.JwtRestController
	pipelineConfigAggregate aggregate.PipelineConfigAggregate
	secretAggregate         aggregate.SecretAggregate
}

func init() {
	app.Register(newPipelineConfigController)
}

func newPipelineConfigController(pipelineConfigAggregate aggregate.PipelineConfigAggregate, secretAggregate aggregate.SecretAggregate) *PipelineConfigController {
	return &PipelineConfigController{
		pipelineConfigAggregate: pipelineConfigAggregate,
		secretAggregate:         secretAggregate,
	}
}

func (c *PipelineConfigController) GetByNameNamespace(name, namespace string) (model.Response, error) {
	response := new(model.BaseResponse)
	pipeline, err := c.pipelineConfigAggregate.Get(name, namespace)
	response.SetData(pipeline)
	return response, err
}

func (c *PipelineConfigController) Post(cmd *command.PipelineStart, properties *jwt.TokenProperties) (response model.Response, err error) {
	log.Debugf("starter pipeline : %v", cmd)
	response = new(model.BaseResponse)
	jwtProps, ok := properties.Items()
	if ok {
		username := jwtProps["username"]
		password := jwtProps["password"]
		token := jwtProps["access_token"]
		secret := &command.Secret{
			Username:  username,
			Password:  password,
			Name:      cmd.Name,
			Namespace: cmd.Namespace,
			Token:     token,
		}
		err = c.secretAggregate.Create(secret)
		if err != nil {
			return
		}
		if len(cmd.Context) == 0 {
			cmd.ParentModule = cmd.Name
			go func() {
				_, err = c.pipelineConfigAggregate.StartPipelineConfig(cmd)
			}()
			return
		}
		for _, ct := range cmd.Context {
			go func() {
				cmd.ParentModule = cmd.Name
				cmd.Path = ct
				paths := strings.Split(ct, "/")
				cmd.Name = paths[len(paths) - 1]
				_, err = c.pipelineConfigAggregate.StartPipelineConfig(cmd)
			}()
		}

	}
	return
}
