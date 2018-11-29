package controller

import (
	"github.com/jinzhu/copier"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/mio/console/pkg/aggregate"
	"hidevops.io/mio/console/pkg/command"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
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
