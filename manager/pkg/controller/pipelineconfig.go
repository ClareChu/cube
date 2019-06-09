package controller

import (
	"github.com/prometheus/common/log"
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/starter/jwt"
)

type PipelineConfigController struct {
	at.JwtRestController
	pipelineConfigAggregate aggregate.PipelineConfigAggregate
	startAggregate          aggregate.StartAggregate
}

func init() {
	app.Register(newPipelineConfigController)
}

func newPipelineConfigController(pipelineConfigAggregate aggregate.PipelineConfigAggregate, startAggregate aggregate.StartAggregate) *PipelineConfigController {
	return &PipelineConfigController{
		pipelineConfigAggregate: pipelineConfigAggregate,
		startAggregate:          startAggregate,
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
		err = c.startAggregate.Init(cmd, jwtProps)
	}
	return
}


func (c *PipelineConfigController) PostRun(cmd *command.PipelineStart, properties *jwt.TokenProperties) (response model.Response, err error) {
	log.Debugf("starter pipeline : %v", cmd)
	response = new(model.BaseResponse)
	jwtProps, ok := properties.Items()
	if ok {
		err = c.startAggregate.Init(cmd, jwtProps)
	}
	return
}