package controller

import (
	"github.com/jinzhu/copier"
	"hidevops.io/cube/console/pkg/aggregate"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type PipelineController struct {
	at.RestController
	pipelineAggregate aggregate.PipelineAggregate
}

func init() {
	app.Register(newPipelineController)
}

func newPipelineController(pipelineAggregate aggregate.PipelineAggregate) *PipelineController {
	return &PipelineController{
		pipelineAggregate: pipelineAggregate,
	}
}

func (p *PipelineController) Post(pipeline *command.PipelineConfigTemplate) (model.Response, error) {
	response := new(model.BaseResponse)
	pic := &v1alpha1.PipelineConfig{}
	copier.Copy(pic, pipeline)
	pc, err := p.pipelineAggregate.Create(pic, "")
	response.SetData(pc)
	return response, err
}
