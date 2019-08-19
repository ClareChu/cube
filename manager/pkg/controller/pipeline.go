package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/starter/jwt"
)

type PipelineController struct {
	at.RestController
	pipelineAggregate aggregate.PipelineAggregate
	startAggregate    aggregate.StartAggregate
}

func init() {
	app.Register(newPipelineController)
}

func newPipelineController(pipelineAggregate aggregate.PipelineAggregate, startAggregate aggregate.StartAggregate) *PipelineController {
	return &PipelineController{
		pipelineAggregate: pipelineAggregate,
		startAggregate:    startAggregate,
	}
}

func (p *PipelineController) Post(cmd *command.PipelineStart, properties *jwt.TokenProperties, ctx context.Context) (response model.Response, err error) {
	log.Debugf("starter pipeline : %v", cmd)
	response = new(model.BaseResponse)
	jwtProps, ok := properties.Items()
	if ok {
		token := ctx.GetHeader(constant.Authorization)
		cmd.Token = token
		err = p.startAggregate.Init(cmd, jwtProps)
	}
	return
}
