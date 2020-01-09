package v2

import (
	"github.com/jinzhu/copier"
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/aggregate/dispatch"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
)

type Selector struct {
	buildConfigAggregate     aggregate.BuildConfigAggregate
	pipelineFactoryInterface dispatch.PipelineFactoryInterface
	pipelineBuilder          builder.PipelineBuilder
}

type SelectorInterface interface {
	Handle(pipeline *v1alpha1.Pipeline) (err error)
}

func init() {
	app.Register(NewSelectorService)
}

func NewSelectorService(buildConfigAggregate aggregate.BuildConfigAggregate,
	pipelineFactoryInterface dispatch.PipelineFactoryInterface) SelectorInterface {
	return &Selector{
		buildConfigAggregate:     buildConfigAggregate,
		pipelineFactoryInterface: pipelineFactoryInterface,
	}
}

func (s *Selector) Handle(pipeline *v1alpha1.Pipeline) (err error) {
	eventType := v1alpha1.Events{}
	if len(pipeline.Status.Stages) == 0 {
		eventType = pipeline.Spec.Events[0]
	} else if pipeline.Status.Phase == constant.Success && len(pipeline.Status.Stages) != len(pipeline.Spec.Events) {
		eventType = pipeline.Spec.Events[len(pipeline.Status.Stages)]
	}
	log.Debugf("EventTypes : %v", eventType.EventTypes)
	params := s.InitReqParams(pipeline, eventType.Name)
	aggregate := s.pipelineFactoryInterface.Get(eventType.EventTypes)
	if aggregate == nil {
		log.Infof("pipeline is complete !!!")
		return
	}
	err = s.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, eventType.EventTypes, constant.Created, "")
	go func() {
		err = aggregate.Create(params)
	}()
	return
}

func (s *Selector) InitReqParams(pipeline *v1alpha1.Pipeline, eventType string) (params *command.PipelineReqParams) {
	params = &command.PipelineReqParams{}
	err := copier.Copy(params, &pipeline.Spec)
	if err != nil {
		log.Errorf("copy is err :%v", err)
		return
	}
	params.EventType = eventType
	params.Name = pipeline.Labels[constant.PipelineConfigName]
	params.PipelineName = pipeline.Name
	params.Namespace = pipeline.Namespace
	params.BuildVersion = pipeline.Labels[constant.BuildPipeline]
	return
}
