package builder

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/openshift"
)

type DeploymentConfigBuilder interface {
	CreateApp(deploy *v1alpha1.Deployment) error
}

func init() {
	app.Register(NewDeploymentConfig)
}

func NewDeploymentConfig(deploymentConfig *openshift.DeploymentConfig, deploymentBuilder DeploymentBuilder) DeploymentConfigBuilder {
	return &DeploymentConfig{
		deploymentConfig:  deploymentConfig,
		deploymentBuilder: deploymentBuilder,
	}
}

type DeploymentConfig struct {
	DeploymentConfigBuilder
	deploymentConfig  *openshift.DeploymentConfig
	deploymentBuilder DeploymentBuilder
}

func (d *DeploymentConfig) CreateApp(deploy *v1alpha1.Deployment) error {
	phase := constant.Success
	app := deploy.Labels[constant.DeploymentConfig]
	fullName := fmt.Sprintf("%s-%s", app, deploy.Spec.Version)
	labels := map[string]string{
		"app":     app,
		"version": deploy.Spec.Version,
	}
	de := &openshift.DeploymentRequest{
		Name:           app,
		Namespace:      deploy.Namespace,
		FullName:       fullName,
		Version:        deploy.Spec.Version,
		Env:            deploy.Spec.Env,
		Labels:         labels,
		Ports:          deploy.Spec.Port,
		Replicas:       1,
		Force:          true,
		HealthEndPoint: "",
		NodeSelector:   "",
		Tag:            deploy.ObjectMeta.Labels[constant.BuildVersion],
	}
	err := d.deploymentConfig.Create(de)
	if err != nil {
		log.Errorf("create openshift deployment config : %v", err)
		phase = constant.Fail
	}
	err = d.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.Deploy, phase)
	return err
}
