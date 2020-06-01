package builder

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	v1 "hidevops.io/cube/manager/pkg/service/apps/v1"
	"hidevops.io/cube/manager/pkg/service/client"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	appsV1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
	podBuilder       PodBuilder
}

func init() {
	app.Register(newDeploymentService)
}

func newDeploymentService(deploymentClient *cube.Deployment, deployment *kube.Deployment, imageStream *cube.ImageStream, podBuilder PodBuilder) DeploymentBuilder {
	return &Deployment{
		deploymentClient: deploymentClient,
		deployment:       deployment,
		imageStream:      imageStream,
		podBuilder:       podBuilder,
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
	err = copier.Copy(dd, &deploy.Spec)
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

func (d *Deployment) Create(dd *command.DeployData) (*appsV1.Deployment, error) {
	runAsRoot := false
	var containers []corev1.Container
	if dd.InitContainer.Name != "" {
		containers = append(containers, dd.InitContainer)
	} else {
		containers = nil
	}
	dpm := &appsV1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", dd.Name, dd.Version),
			Namespace: dd.Namespace,
			Labels: map[string]string{
				"app":     dd.Name,
				"version": dd.Version,
			},
		},

		Spec: appsV1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     dd.Name,
					"version": dd.Version,
				},
			},
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
	return d.CreateDeployment(dd, dpm)
}

func (d *Deployment) CreateDeployment(dd *command.DeployData, dpm *appsV1.Deployment) (*appsV1.Deployment, error) {
	clientSet, err := client.GetDefaultK8sClientSet()
	deployment := v1.Deployment{
		ClientSet: clientSet,
	}
	dp, err := deployment.Get(fmt.Sprintf("%s-%s", dd.Name, dd.Version), dd.Namespace, metav1.GetOptions{})
	log.Infof("*** deploy is exist *** %s", dd.ForceUpdate)
	for i := 0; i < 3; i++ {
		if err == nil {
			dpm.ObjectMeta = dp.ObjectMeta
			log.Infof("****  update deploy app deployment  ***")

			e := deployment.Update(dpm)
			if e != nil {
				log.Errorf("*** update deploy error: %v try again ***", e)
				time.Sleep(10 * time.Second)
				continue
			}
			return dp, nil

		} else {
			log.Infof("****  create deploy app deployment  ***")
			dp, e := deployment.Create(dpm)
			if e != nil {
				log.Errorf("*** create deploy error: %v try again ***", e)
				time.Sleep(10 * time.Second)
				continue
			}
			return dp, nil
		}
	}
	return dp, errors.New("create deployment and update deployment is error")
}