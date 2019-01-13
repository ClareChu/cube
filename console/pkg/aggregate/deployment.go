package aggregate

import (
	"errors"
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/copier"
	"hidevops.io/mio/console/pkg/builder"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"os"
	"strconv"
	"time"
)

type DeploymentAggregate interface {
	Create(deploymentConfig *v1alpha1.DeploymentConfig, pipelineName, version, buildVersion string) (deployment *v1alpha1.Deployment, err error)
	Watch(name, namespace string) error
	Selector(deploy *v1alpha1.Deployment) error
	CreateDeployment(deploy *v1alpha1.Deployment) (err error)
}

type Deployment struct {
	DeploymentAggregate
	deploymentClient        *mio.Deployment
	remoteAggregate         RemoteAggregate
	pipelineBuilder         builder.PipelineBuilder
	deploymentBuilder       builder.DeploymentBuilder
	deploymentConfigBuilder builder.DeploymentConfigBuilder
	tagAggregate            TagAggregate
}

func init() {
	app.Register(NewDeploymentService)
}

func NewDeploymentService(deploymentClient *mio.Deployment, remoteAggregate RemoteAggregate, deploymentBuilder builder.DeploymentBuilder, pipelineBuilder builder.PipelineBuilder, deploymentConfigBuilder builder.DeploymentConfigBuilder, tagAggregate TagAggregate) DeploymentAggregate {
	return &Deployment{
		deploymentClient:        deploymentClient,
		remoteAggregate:         remoteAggregate,
		deploymentBuilder:       deploymentBuilder,
		pipelineBuilder:         pipelineBuilder,
		deploymentConfigBuilder: deploymentConfigBuilder,
		tagAggregate:            tagAggregate,
	}
}

func (d *Deployment) Create(deploymentConfig *v1alpha1.DeploymentConfig, pipelineName, version, buildVersion string) (deployment *v1alpha1.Deployment, err error) {
	log.Debugf("deployment config create :%v", deploymentConfig)
	number := fmt.Sprintf("%d", deploymentConfig.Status.LastVersion)
	nameVersion := fmt.Sprintf("%s-%s", deploymentConfig.Name, number)
	deployment = new(v1alpha1.Deployment)
	copier.Copy(deployment, deploymentConfig)
	deployment.TypeMeta = v1.TypeMeta{
		Kind:       constant.BuildKind,
		APIVersion: constant.BuildApiVersion,
	}
	deployment.ObjectMeta = v1.ObjectMeta{
		Name:      nameVersion,
		Namespace: deploymentConfig.Namespace,
		Labels: map[string]string{
			constant.App:              nameVersion,
			constant.Number:           number,
			constant.DeploymentConfig: deploymentConfig.Name,
			constant.PipelineName:     pipelineName,
			constant.Version:          version,
			constant.BuildVersion:     buildVersion,
		},
	}
	deployment.Spec.Version = version
	deployment.Spec.Tag = version
	config, err := d.deploymentClient.Create(deployment)
	if err != nil {
		log.Errorf("create build error :%v", err)
		return
	} else {
		err = d.Watch(config.Name, config.Namespace)
	}
	return config, err
}

func (d *Deployment) Watch(name, namespace string) error {
	log.Debugf("build config Watch :%v", name)
	kubeWatchTimeout, err := strconv.Atoi(os.Getenv("KUBE_WATCH_TIMEOUT"))
	after := time.Duration(kubeWatchTimeout) * time.Minute

	timeout := int64(constant.TimeoutSeconds)
	option := v1.ListOptions{
		TimeoutSeconds: &timeout,
		LabelSelector:  fmt.Sprintf("app=%s", name),
	}
	w, err := d.deploymentClient.Watch(option, namespace)
	if err != nil {
		return nil
	}
	for {
		select {
		case <-time.After(after):
			return errors.New("pod query timeout 10 minutes")
		case event, ok := <-w.ResultChan():
			if !ok {
				log.Info(" build watch resultChan: ", ok)
				return nil
			}
			switch event.Type {
			case watch.Added:
				deploy := event.Object.(*v1alpha1.Deployment)
				err = d.Selector(deploy)
				if err != nil {
					log.Errorf("add selector : %v", err)
					return err
				}
				log.Infof("event type :%v, err: %v", deploy.Status, err)
			case watch.Modified:
				deploy := event.Object.(*v1alpha1.Deployment)
				if deploy.Status.Phase == constant.Fail {
					return fmt.Errorf("build status phase error: %s", deploy.Status.Phase)
				}
				err = d.Selector(deploy)
				if err != nil {
					log.Errorf("add selector : %v", err)
					return err
				}
				log.Infof("event type :%v", deploy.Status)
			case watch.Deleted:
				log.Info("Deleted: ", event.Object)
				return nil
			default:
				log.Error("Failed")
			}
		}
	}
}

func (d *Deployment) Selector(deploy *v1alpha1.Deployment) error {
	eventType := "Default"
	if len(deploy.Status.Stages) == 0 {
		if len(deploy.Spec.EnvType) == 0 {
			return fmt.Errorf("not fount events")
		}
		eventType = deploy.Spec.EnvType[0]
	} else if deploy.Status.Phase == constant.Success && len(deploy.Status.Stages) != len(deploy.Spec.EnvType) {
		eventType = deploy.Spec.EnvType[len(deploy.Status.Stages)]
	} else if len(deploy.Status.Stages) == len(deploy.Spec.EnvType) && deploy.Status.Phase == constant.Success {
		eventType = constant.Ending
	}
	var err error
	switch eventType {
	case constant.RemoteDeploy:
		go func() {
			d.tagAggregate.TagImage(deploy)
		}()
		err = d.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.RemoteDeploy, constant.Created)
	case constant.Deploy:
		go func() {
			d.CreateDeployment(deploy)
		}()
		err = d.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.Deploy, constant.Created)
	case constant.Ending:
		err = d.pipelineBuilder.Update(deploy.ObjectMeta.Labels[constant.PipelineName], deploy.Namespace, constant.Deploy, deploy.Status.Phase, deploy.ObjectMeta.Labels[constant.Number])
		err = fmt.Errorf("build is ending")
	default:

	}
	return err
}

func (d *Deployment) CreateDeployment(deploy *v1alpha1.Deployment) (err error) {
	if os.Getenv("IS_OPENSHIFT") == "" {
		err = d.deploymentBuilder.CreateApp(deploy)
		return
	}
	err = d.deploymentConfigBuilder.CreateApp(deploy)
	return
}
