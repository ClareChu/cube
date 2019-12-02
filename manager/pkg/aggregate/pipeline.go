package aggregate

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
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
	InitReqParams(pipeline *v1alpha1.Pipeline, eventType string) (params *command.PipelineReqParams)
}

type Pipeline struct {
	PipelineAggregate
	pipelineClient    *cube.Pipeline
	selectorInterface SelectorInterface
}

func init() {
	app.Register(NewPipelineService)
}

func NewPipelineService(pipelineClient *cube.Pipeline,
	selectorInterface SelectorInterface) PipelineAggregate {
	return &Pipeline{
		pipelineClient:    pipelineClient,
		selectorInterface: selectorInterface,
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
	err := copier.Copy(&pipeline, pipelineConfig)
	if err != nil {
		log.Errorf("Copy pipeline is err : %v", err)
		return nil, err
	}
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
				err = p.selectorInterface.Handle(pipeline)
				if err != nil {
					return
				}
			case watch.Modified:
				pipeline = event.Object.(*v1alpha1.Pipeline)
				log.Infof("Modified event type :%v", pipeline.Status)
				err = p.selectorInterface.Handle(pipeline)
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
