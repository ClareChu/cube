package builder

import (
	"fmt"
	"github.com/prometheus/common/log"
	"hidevops.io/cube/console/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
	"time"
)

type DeploymentBuilder interface {
	Update(name, namespace, event, phase string) error
	CreateApp(deploy *v1alpha1.Deployment) error
}

type Deployment struct {
	DeploymentBuilder
	deploymentClient *cube.Deployment
	deployment       *kube.Deployment
	imageStream      *cube.ImageStream
}

func init() {
	app.Register(newDeploymentService)
}

func newDeploymentService(deploymentClient *cube.Deployment, deployment *kube.Deployment, imageStream *cube.ImageStream) DeploymentBuilder {
	return &Deployment{
		deploymentClient: deploymentClient,
		deployment:       deployment,
		imageStream:      imageStream,
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
	namespace := GetNamespace(deploy.Namespace, deploy.Spec.Profile)
	i, err := d.imageStream.Get(deploy.Labels[constant.DeploymentConfig], namespace)
	if err != nil {
		log.Errorf("get image err: %v", err)
		return err
	}
	request := &kube.DeployRequest{
		App:            deploy.Labels[constant.DeploymentConfig],
		Namespace:      namespace,
		Ports:          deploy.Spec.Port,
		Replicas:       deploy.Spec.Replicas,
		Version:        deploy.Spec.Version,
		Labels:         deploy.Spec.Labels,
		ReadinessProbe: deploy.Spec.ReadinessProbe,
		NodeSelector:   deploy.Spec.NodeSelector,
		LivenessProbe:  deploy.Spec.LivenessProbe,
		Env:            deploy.Spec.Env,
		Volumes:        deploy.Spec.Volumes,
		VolumeMounts:   deploy.Spec.VolumeMounts,
		DockerImage:    i.Spec.Tags[constant.Latest].DockerImageReference,
	}
	_, err = d.deployment.Deploy(request)
	if err != nil {
		log.Errorf("create app :%v", err)
		phase = constant.Fail
	}
	err = d.Update(deploy.Name, deploy.Namespace, constant.Deploy, phase)
	log.Debugf("create app update pipeline :name %v,namespace %v,deploy %v, type:%v, error %v", deploy.Name, deploy.Namespace, constant.Deploy, phase, err)
	return err
}

func GetNamespace(space, profile string) string {
	if profile == "" {
		return space
	}
	return fmt.Sprintf("%s-%s", space, profile)
}
