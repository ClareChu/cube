package aggregate

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"os"
	"strconv"
	"time"
)

type PipelineAggregate interface {
	Get(name, namespace string) (*v1alpha1.Pipeline, error)
	Watch(name, namespace string) (pipeline *v1alpha1.Pipeline, err error)
	Create(pipelineConfig *v1alpha1.PipelineConfig, templateName string) (*v1alpha1.Pipeline, error)
	Selector(pipeline *v1alpha1.Pipeline) (err error)
	InitReqParams(pipeline *v1alpha1.Pipeline, eventType string) (params *command.PipelineReqParams)
}

type Pipeline struct {
	PipelineAggregate
	pipelineClient            *cube.Pipeline
	buildConfigAggregate      BuildConfigAggregate
	deploymentConfigAggregate DeploymentConfigAggregate
	pipelineBuilder           builder.PipelineBuilder
	serviceConfigAggregate    ServiceConfigAggregate
	gatewayConfigAggregate    GatewayConfigAggregate
	imageStreamAggregate      ImageStreamAggregate
	volumeAggregate           VolumeAggregate
	callbackAggregate         CallbackAggregate
}

func init() {
	app.Register(NewPipelineService)
}

func NewPipelineService(pipelineClient *cube.Pipeline,
	buildConfigAggregate BuildConfigAggregate, deploymentConfigAggregate DeploymentConfigAggregate,
	pipelineBuilder builder.PipelineBuilder, serviceConfigAggregate ServiceConfigAggregate,
	gatewayConfigAggregate GatewayConfigAggregate, imageStreamAggregate ImageStreamAggregate,
	volumeAggregate VolumeAggregate,
	callbackAggregate CallbackAggregate) PipelineAggregate {
	return &Pipeline{
		pipelineClient:            pipelineClient,
		buildConfigAggregate:      buildConfigAggregate,
		deploymentConfigAggregate: deploymentConfigAggregate,
		pipelineBuilder:           pipelineBuilder,
		serviceConfigAggregate:    serviceConfigAggregate,
		gatewayConfigAggregate:    gatewayConfigAggregate,
		imageStreamAggregate:      imageStreamAggregate,
		volumeAggregate:           volumeAggregate,
		callbackAggregate:         callbackAggregate,
	}
}

func (p *Pipeline) Get(name, namespace string) (*v1alpha1.Pipeline, error) {
	log.Debug("build config create :%v", name, namespace)
	config, err := p.pipelineClient.Get(name, namespace)
	return config, err
}

func (p *Pipeline) Create(pipelineConfig *v1alpha1.PipelineConfig, templateName string) (*v1alpha1.Pipeline, error) {
	log.Debugf("create pipeline :%v", pipelineConfig)
	number := fmt.Sprintf("%d", pipelineConfig.Status.LastVersion)
	nameVersion := pipelineConfig.Name + "-" + number
	pipeline := new(v1alpha1.Pipeline)
	copier.Copy(&pipeline, pipelineConfig)
	pipeline.TypeMeta = v1.TypeMeta{
		Kind:       constant.PipelineKind,
		APIVersion: constant.PipelineApiVersion,
	}
	pipeline.ObjectMeta = v1.ObjectMeta{
		Name:      nameVersion,
		Namespace: pipelineConfig.Namespace,
		Labels: map[string]string{
			constant.App:                nameVersion,
			constant.Version:            pipelineConfig.Spec.Version,
			constant.Number:             number,
			constant.PipelineConfigName: pipelineConfig.Name,
			constant.CodeType:           templateName,
		},
	}
	pipeline.Status = v1alpha1.PipelineStatus{
		Name:      pipelineConfig.Name,
		Namespace: pipelineConfig.Namespace,
	}
	config, err := p.pipelineClient.Create(pipeline)
	if err != nil {
		//TODO 启动 pipeline watch
		log.Errorf("create pipeline error :%v", err)
	}
	_, err = p.Watch(nameVersion, pipelineConfig.Namespace)
	return config, err
}

func (p *Pipeline) Watch(name, namespace string) (pipeline *v1alpha1.Pipeline, err error) {
	kubeWatchTimeout, err := strconv.Atoi(os.Getenv("KUBE_WATCH_TIMEOUT"))
	after := time.Duration(kubeWatchTimeout) * time.Minute

	timeout := int64(constant.TimeoutSeconds)
	options := v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:  fmt.Sprintf("app=%s", name),
	}
	w, err := p.pipelineClient.Watch(options, namespace)
	if err != nil {
		return
	}
	for {
		select {
		case <-time.After(after):
			return nil, errors.New("10 min")
		case event, ok := <-w.ResultChan():
			if !ok {
				log.Info("pipeline resultChan: ", ok)
				return nil, nil
			}
			switch event.Type {
			case watch.Added:
				pipeline = event.Object.(*v1alpha1.Pipeline)
				log.Infof("add event type :%v", pipeline.Status)
				err = p.Selector(pipeline)
				if err != nil {
					return
				}
			case watch.Modified:
				pipeline = event.Object.(*v1alpha1.Pipeline)
				log.Infof("Modified event type :%v", pipeline.Status)
				err = p.Selector(pipeline)
				if err != nil {
					return
				}
			case watch.Deleted:
				log.Info("Deleted: ", event.Object)
				return
			default:
				log.Error("Failed")
				return
			}
		}
	}
}

func (p *Pipeline) Selector(pipeline *v1alpha1.Pipeline) (err error) {
	log.Infof("pipeline selector : %v", pipeline)
	eventType := v1alpha1.Events{}
	if len(pipeline.Status.Stages) == 0 {
		eventType = pipeline.Spec.Events[0]
	} else if pipeline.Status.Phase == constant.Success && len(pipeline.Status.Stages) != len(pipeline.Spec.Events) {
		eventType = pipeline.Spec.Events[len(pipeline.Status.Stages)]
	}
	log.Debugf("EventTypes : %v", eventType.EventTypes)
	params := p.InitReqParams(pipeline, eventType.Name)
	switch eventType.EventTypes {
	case constant.BuildPipeline:
		go func() {
			p.buildConfigAggregate.Create(params)
		}()
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.BuildPipeline, constant.Created, "")
		return
	case constant.Deploy:
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Deploy, constant.Created, "")
		go func() {
			p.deploymentConfigAggregate.Create(params, pipeline.Labels[constant.BuildPipeline])
		}()
		return
	case constant.Service:
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Service, constant.Created, "")
		go func() {
			p.serviceConfigAggregate.Create(params)
		}()
		return
	case constant.Gateway:
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Gateway, constant.Created, "")
		go func() {
			p.gatewayConfigAggregate.Create(params)
		}()

	case constant.ImageStream:
		go func() {
			p.imageStreamAggregate.CreateImage(params)
		}()
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.ImageStream, constant.Created, "")
	case constant.Volume:
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Volume, constant.Created, "")
		go func() {
			p.volumeAggregate.Create(params)
		}()
	case constant.CallBack:
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.CallBack, constant.Created, "")
		go func() {
			p.callbackAggregate.Create(params)
		}()
	default:
		return
	}
	return
}

func (p *Pipeline) InitReqParams(pipeline *v1alpha1.Pipeline, eventType string) (params *command.PipelineReqParams) {
	params = &command.PipelineReqParams{}
	copier.Copy(params, &pipeline.Spec)
	params.EventType = eventType
	params.Name = pipeline.Labels[constant.PipelineConfigName]
	params.PipelineName = pipeline.Name
	params.Namespace = pipeline.Namespace
	return
}
