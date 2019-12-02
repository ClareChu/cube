package aggregate

import (
	"github.com/jinzhu/copier"
	"hidevops.io/cube/manager/pkg/aggregate/dispatch"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
)

type Selector struct {
	buildConfigAggregate      BuildConfigAggregate
	deploymentConfigAggregate DeploymentConfigAggregate
	serviceConfigAggregate    ServiceConfigAggregate
	gatewayConfigAggregate    GatewayConfigAggregate
	imageStreamAggregate      ImageStreamAggregate
	volumeAggregate           VolumeAggregate
	callbackAggregate         CallbackAggregate
	pipelineBuilder           builder.PipelineBuilder
	pipelineFactoryInterface  dispatch.PipelineFactoryInterface
}

type SelectorInterface interface {
	Handle(pipeline *v1alpha1.Pipeline) (err error)
}

func init() {
	app.Register(NewSelectorService)
}

func NewSelectorService(buildConfigAggregate BuildConfigAggregate,
	deploymentConfigAggregate DeploymentConfigAggregate,
	serviceConfigAggregate ServiceConfigAggregate,
	gatewayConfigAggregate GatewayConfigAggregate,
	imageStreamAggregate ImageStreamAggregate,
	volumeAggregate VolumeAggregate,
	callbackAggregate CallbackAggregate,
	pipelineBuilder builder.PipelineBuilder,
	pipelineFactoryInterface dispatch.PipelineFactoryInterface) SelectorInterface {
	return &Selector{
		buildConfigAggregate:      buildConfigAggregate,
		deploymentConfigAggregate: deploymentConfigAggregate,
		serviceConfigAggregate:    serviceConfigAggregate,
		gatewayConfigAggregate:    gatewayConfigAggregate,
		imageStreamAggregate:      imageStreamAggregate,
		volumeAggregate:           volumeAggregate,
		callbackAggregate:         callbackAggregate,
		pipelineBuilder:           pipelineBuilder,
		pipelineFactoryInterface:  pipelineFactoryInterface,
	}
}

func (s *Selector) HandleV2(pipeline *v1alpha1.Pipeline) (err error) {
	log.Infof("pipeline selector : %v", pipeline)
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

func (s *Selector) Handle(pipeline *v1alpha1.Pipeline) (err error) {
	log.Infof("pipeline selector : %v", pipeline)
	eventType := v1alpha1.Events{}
	if len(pipeline.Status.Stages) == 0 {
		eventType = pipeline.Spec.Events[0]
	} else if pipeline.Status.Phase == constant.Success && len(pipeline.Status.Stages) != len(pipeline.Spec.Events) {
		eventType = pipeline.Spec.Events[len(pipeline.Status.Stages)]
	}
	log.Debugf("EventTypes : %v", eventType.EventTypes)
	params := s.InitReqParams(pipeline, eventType.Name)
	switch eventType.EventTypes {

	case constant.BuildPipeline:
		err = s.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.BuildPipeline, constant.Created, "")
		go func() {
			err = s.buildConfigAggregate.Create(params)
		}()
		return
	case constant.Deploy:
		err = s.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Deploy, constant.Created, "")
		go func() {
			err = s.deploymentConfigAggregate.Create(params)
		}()
		return
	case constant.Service:
		err = s.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Service, constant.Created, "")
		go func() {
			err = s.serviceConfigAggregate.Create(params)
		}()
		return
	case constant.Gateway:
		err = s.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Gateway, constant.Created, "")
		go func() {
			err = s.gatewayConfigAggregate.Create(params)

		}()

	case constant.ImageStream:
		err = s.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.ImageStream, constant.Created, "")
		go func() {
			err = s.imageStreamAggregate.Create(params)
		}()
	case constant.Volume:
		err = s.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Volume, constant.Created, "")
		go func() {
			err = s.volumeAggregate.Create(params)
		}()
	case constant.CallBack:
		err = s.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.CallBack, constant.Created, "")
		go func() {
			err = s.callbackAggregate.Create(params)
		}()
	default:
		return
	}
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
