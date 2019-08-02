package builder

import (
	"fmt"
	"github.com/jinzhu/copier"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	extensionsV1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	i, err := d.imageStream.Get(deploy.Labels[constant.DeploymentConfig], deploy.Namespace)
	if err != nil {
		log.Errorf("get image err: %v", err)
		return err
	}
	deploy.Spec.Container.Image = i.Spec.Tags[constant.Latest].DockerImageReference
	dd := &command.DeployData{}
	copier.Copy(dd, &deploy.Spec)
	dd.Name = deploy.Labels[constant.DeploymentConfig]
	dd.Namespace = deploy.Namespace
	_, err = d.Create(dd)
	if err != nil {
		log.Errorf("create app :%v", err)
		phase = constant.Fail
	}
	err = d.Update(deploy.Name, deploy.Namespace, constant.Deploy, phase)
	log.Debugf("create app update pipeline :name %v,namespace %v,deploy %v, type:%v, error %v", deploy.Name, deploy.Namespace, constant.Deploy, phase, err)
	return err
}

func int32Ptr(i int32) *int32 { return &i }

func (d *Deployment) Create(dd *command.DeployData) (*extensionsV1beta1.Deployment, error) {
	runAsRoot := false
	containers := []corev1.Container{}
	if dd.InitContainer.Name != "" {
		containers = append(containers, dd.InitContainer)
	} else {
		containers = nil
	}
	dpm := &extensionsV1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", dd.Name, dd.Version),
			Namespace: dd.Namespace,
			Labels: map[string]string{
				"app":     dd.Name,
				"version": dd.Version,
			},
		},

		Spec: extensionsV1beta1.DeploymentSpec{
			Replicas:             dd.Replicas,
			Strategy:             dd.Strategy,
			RevisionHistoryLimit: int32Ptr(10),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: dd.Name,
					Labels: map[string]string{
						"app":     dd.Name,
						"version": dd.Version,
					},
				},
				Spec: corev1.PodSpec{
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: &runAsRoot,
					},
					NodeSelector: dd.NodeSelector,
					Containers: []corev1.Container{
						dd.Container,
					},
					InitContainers: containers,
					Volumes:        dd.Volumes,
				},
			},
		},
	}

	dp, err := d.deployment.Get(fmt.Sprintf("%s-%s", dd.Name, dd.Version), dd.Namespace, metav1.GetOptions{})
	log.Infof("*** deploy is exist *** %s", dd.ForceUpdate)
	if err == nil {
		if dd.ForceUpdate {
			log.Infof("****  update deploy app dpm  ***")
			dpm.ObjectMeta = dp.ObjectMeta
			err = d.deployment.Update(dpm)
			return nil, err
		}
		return dp, err

	}
	log.Infof("****  create deploy app dpm  ***")
	return d.deployment.Create(dpm)
}
