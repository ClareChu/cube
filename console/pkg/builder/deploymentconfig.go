package builder

import (
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/openshift"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
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
	namespace := ""
	if deploy.Spec.Profile == "" {
		namespace = deploy.Namespace
	} else {
		namespace = deploy.Namespace + "-" + deploy.Spec.Profile
	}
	labels := map[string]string{
		"app":     app,
		"version": deploy.Spec.Version,
	}
	de := &openshift.DeploymentRequest{
		Name:           app,
		Namespace:      namespace,
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
		log.Error("create openshift deployment config : %v", err)
		phase = constant.Fail
	}
	err = d.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.Deploy, phase)
	return err
}
