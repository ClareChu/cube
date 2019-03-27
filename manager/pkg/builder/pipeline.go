package builder

import (
	"github.com/prometheus/common/log"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"time"
)

type PipelineBuilder interface {
	Update(name, namespace, eventType, phase, eventVersion string) error
}

type Pipeline struct {
	PipelineBuilder
	pipelineClient *cube.Pipeline
}

func init() {
	app.Register(newPipelineService)
}

func newPipelineService(pipelineClient *cube.Pipeline) PipelineBuilder {
	return &Pipeline{
		pipelineClient: pipelineClient,
	}
}

func (p *Pipeline) Update(name, namespace, eventType, phase, eventVersion string) error {
	log.Debugf("name: %v, namespace: %v , eventType: %v , phase : %v ", name, namespace, eventType, phase)
	pipeline, err := p.pipelineClient.Get(name, namespace)
	if err != nil {
		log.Errorf("get pipeline err : %v", err)
		return err
	}
	if eventVersion != "" {
		pipeline.ObjectMeta.Labels[eventType] = eventVersion
	}
	stage := v1alpha1.Stages{}
	if pipeline.Status.Phase == constant.Created {
		stage = pipeline.Status.Stages[len(pipeline.Status.Stages)-1]
		stage.DurationMilliseconds = time.Now().Unix() - stage.StartTime
		pipeline.Status.Stages[len(pipeline.Status.Stages)-1] = stage
	} else {
		stage = v1alpha1.Stages{
			Name:                 eventType,
			StartTime:            time.Now().Unix(),
			DurationMilliseconds: 0,
		}
		pipeline.Status.Stages = append(pipeline.Status.Stages, stage)
	}
	pipeline.Status.Phase = phase
	_, err = p.pipelineClient.Update(pipeline.Name, pipeline.Namespace, pipeline)
	return err
}
