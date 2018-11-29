package aggregate

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/mio/console/pkg/builder"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"time"
)

type PipelineAggregate interface {
	Get(name, namespace string) (*v1alpha1.Pipeline, error)
	Watch(name, namespace string) (pipeline *v1alpha1.Pipeline, err error)
	Create(pipelineConfig *v1alpha1.PipelineConfig, sourceCode string) (*v1alpha1.Pipeline, error)
	Selector(pipeline *v1alpha1.Pipeline) (err error)
}

type Pipeline struct {
	PipelineAggregate
	pipelineClient            *mio.Pipeline
	buildConfigAggregate      BuildConfigAggregate
	deploymentConfigAggregate DeploymentConfigAggregate
	pipelineBuilder           builder.PipelineBuilder
	serviceConfigAggregate    ServiceConfigAggregate
	gatewayConfigAggregate    GatewayConfigAggregate
}

func init() {
	app.Register(NewPipelineService)
}

func NewPipelineService(pipelineClient *mio.Pipeline, buildConfigAggregate BuildConfigAggregate, deploymentConfigAggregate DeploymentConfigAggregate, pipelineBuilder builder.PipelineBuilder, serviceConfigAggregate ServiceConfigAggregate, gatewayConfigAggregate GatewayConfigAggregate) PipelineAggregate {
	return &Pipeline{
		pipelineClient:            pipelineClient,
		buildConfigAggregate:      buildConfigAggregate,
		deploymentConfigAggregate: deploymentConfigAggregate,
		pipelineBuilder:           pipelineBuilder,
		serviceConfigAggregate:    serviceConfigAggregate,
		gatewayConfigAggregate:    gatewayConfigAggregate,
	}
}

func (p *Pipeline) Get(name, namespace string) (*v1alpha1.Pipeline, error) {
	log.Debug("build config create :%v", name, namespace)
	config, err := p.pipelineClient.Get(name, namespace)
	return config, err
}

func (p *Pipeline) Create(pipelineConfig *v1alpha1.PipelineConfig, sourceCode string) (*v1alpha1.Pipeline, error) {
	log.Debug("create pipeline :%v", pipelineConfig)
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
			constant.CodeType:           sourceCode,
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
		case <-time.After(10 * time.Second):
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
				log.Infof("Deleted: ", event.Object)
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
	log.Debugf("events : %v, eventType name: %v", pipeline.Spec.Events, eventType.Name)
	switch eventType.EventTypes {
	case constant.BuildPipeline:
		go func() {
			p.buildConfigAggregate.Create(pipeline.Labels[constant.PipelineConfigName], pipeline.Name, pipeline.Namespace, eventType.Name, pipeline.Spec.Version)
		}()
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.BuildPipeline, constant.Created, "")
		return
	case constant.Deploy:
		go func() {
			p.deploymentConfigAggregate.Create(pipeline.Labels[constant.PipelineConfigName], pipeline.Name, pipeline.Namespace, eventType.Name, pipeline.Spec.Version, pipeline.Labels[constant.BuildPipeline], pipeline.Spec.Profile)
		}()
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Deploy, constant.Created, "")
		return
	case constant.Service:
		go func() {
			p.serviceConfigAggregate.Create(pipeline.Labels[constant.PipelineConfigName], pipeline.Name, pipeline.Namespace, eventType.Name, pipeline.Spec.Version, pipeline.Spec.Profile)
		}()
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Service, constant.Created, "")
		return
	case constant.Gateway:
		go func() {
			p.gatewayConfigAggregate.Create(pipeline.Labels[constant.PipelineConfigName], pipeline.Name, pipeline.Namespace, eventType.Name, pipeline.Spec.Version, pipeline.Spec.Profile)
		}()
		err = p.pipelineBuilder.Update(pipeline.Name, pipeline.Namespace, constant.Gateway, constant.Created, "")
	default:

		return
	}
	return
}
