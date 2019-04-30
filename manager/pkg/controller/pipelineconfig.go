package controller

import (
	"fmt"
	"github.com/prometheus/common/log"
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
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
	namespaceAggregate      aggregate.NamespaceAggregate
	roleAggregate           aggregate.RoleAggregate
	roleBindingAggregate    aggregate.RoleBindingAggregate
}

func init() {
	app.Register(newPipelineConfigController)
}

func newPipelineConfigController(pipelineConfigAggregate aggregate.PipelineConfigAggregate, secretAggregate aggregate.SecretAggregate, namespaceAggregate aggregate.NamespaceAggregate, roleAggregate aggregate.RoleAggregate, roleBindingAggregate aggregate.RoleBindingAggregate) *PipelineConfigController {
	return &PipelineConfigController{
		pipelineConfigAggregate: pipelineConfigAggregate,
		secretAggregate:         secretAggregate,
		namespaceAggregate:      namespaceAggregate,
		roleAggregate:           roleAggregate,
		roleBindingAggregate:    roleBindingAggregate,
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
		if cmd.Namespace == "" {
			if cmd.Profile == "" {
				cmd.Namespace = cmd.Project
			} else {
				cmd.Namespace = fmt.Sprintf("%s-%s", cmd.Project, cmd.Profile)
			}
		}
		//TODO CREATE NAMESPACE
		err = c.namespaceAggregate.InitNamespace(cmd.Namespace)
		if err != nil {
			return
		}
		//TODO create role
		err = c.roleAggregate.Create(constant.Default ,cmd.Namespace)
		if err != nil {
			return
		}
		//TODO create rolebinding
		err = c.roleBindingAggregate.Create(constant.Default ,cmd.Namespace)
		if err != nil {
			return
		}
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
		if len(cmd.Context) == 0 || cmd.Context[0] == "" {
			cmd.AppRoot = cmd.Name
			go func() {
				_, err = c.pipelineConfigAggregate.StartPipelineConfig(cmd)
			}()
			return
		}
		for _, ct := range cmd.Context {
			log.Info(ct)
			paths := strings.Split(ct, "/")
			name := paths[len(paths)-1]
			command := command.PipelineStart{
				Name:       name,
				Namespace:  cmd.Namespace,
				Version:    cmd.Version,
				Profile:    cmd.Profile,
				Path:       ct,
				AppRoot:    cmd.Name,
				SourceCode: cmd.SourceCode,
				Branch:     cmd.Branch,
				Project:    cmd.Project,
			}
			go func() {
				_, err = c.pipelineConfigAggregate.StartPipelineConfig(&command)
			}()
		}

	}
	return
}
