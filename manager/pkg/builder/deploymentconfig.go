package builder

import (
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
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
	return nil
}
