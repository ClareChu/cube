package builder

import (
	"github.com/prometheus/common/log"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/starter/mio"
	"time"
)

type DeploymentBuilder interface {
	Update(name, namespace, event, phase string) error
}

type Deployment struct {
	DeploymentBuilder
	deploymentClient *mio.Deployment
}

func init() {
	app.Register(newDeploymentService)
}

func newDeploymentService(deploymentClient *mio.Deployment) DeploymentBuilder {
	return &Deployment{
		deploymentClient: deploymentClient,
	}
}

func (d *Deployment) Update(name, namespace, event, phase string) error {
	deploy, err := d.deploymentClient.Get(name, namespace)
	if err != nil {
		return err
	}
	stage := v1alpha1.Stages{}
	if deploy.Status.Phase == constant.Created {
		stage = deploy.Status.Stages[len(deploy.Status.Stages)-1]
		stage.DurationMilliseconds = time.Now().Unix() - stage.StartTime
		deploy.Status.Stages[len(deploy.Status.Stages)-1] = stage
	} else {
		stage = v1alpha1.Stages{
			Name:                 event,
			StartTime:            time.Now().Unix(),
			DurationMilliseconds: 0,
		}
		deploy.Status.Stages = append(deploy.Status.Stages, stage)
	}
	deploy.Status.Phase = phase
	_, err = d.deploymentClient.Update(name, namespace, deploy)
	if err != nil {
		log.Errorf("deployment update err :%v", err)
	}
	return err
}
