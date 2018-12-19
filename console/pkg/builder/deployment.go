package builder

import (
	"github.com/prometheus/common/log"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/starter/mio"
	"time"
)

type DeploymentBuilder interface {
	Update(name, namespace, event, phase string) error
	CreateApp(deploy *v1alpha1.Deployment) error
}

type Deployment struct {
	DeploymentBuilder
	deploymentClient *mio.Deployment
	deployment       *kube.Deployment
}

func init() {
	app.Register(newDeploymentService)
}

func newDeploymentService(deploymentClient *mio.Deployment, deployment *kube.Deployment) DeploymentBuilder {
	return &Deployment{
		deploymentClient: deploymentClient,
		deployment:       deployment,
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

func (d *Deployment) CreateApp(deploy *v1alpha1.Deployment) error {
	phase := constant.Success
	namespace := ""
	if deploy.Spec.Profile == "" {
		namespace = deploy.Namespace
	} else {
		namespace = deploy.Namespace + "-" + deploy.Spec.Profile
	}
	request := &kube.DeployRequest{
		App:            deploy.Labels[constant.DeploymentConfig],
		Namespace:      namespace,
		Ports:          deploy.Spec.Port,
		Replicas:       deploy.Spec.Replicas,
		Version:        deploy.Spec.Version,
		Tag:            deploy.ObjectMeta.Labels[constant.BuildVersion],
		Labels:         deploy.Spec.Labels,
		ReadinessProbe: deploy.Spec.ReadinessProbe,
		NodeSelector:   deploy.Spec.NodeSelector,
		LivenessProbe:  deploy.Spec.LivenessProbe,
		Env:            deploy.Spec.Env,
		DockerRegistry: deploy.Spec.DockerRegistry,
		Volumes:        deploy.Spec.Volumes,
		VolumeMounts:   deploy.Spec.VolumeMounts,
	}
	_, err := d.deployment.Deploy(request)
	if err != nil {
		log.Errorf("create app :%v", err)
		phase = constant.Fail
	}
	err = d.Update(deploy.Name, deploy.Namespace, constant.Deploy, phase)
	log.Debugf("create app update pipeline :name %v,namespace %v,deploy %v, type:%v, error %v", deploy.Name, deploy.Namespace, constant.Deploy, phase, err)
	return err
}
