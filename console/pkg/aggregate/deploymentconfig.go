package aggregate

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/copier"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/console/pkg/builder"
	"hidevops.io/mio/console/pkg/command"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentConfigAggregate interface {
	Template(cmd *command.DeploymentConfig) (deploymentConfig *v1alpha1.DeploymentConfig, err error)
	Create(name, pipelineName, namespace, sourceType, version, buildVersion, profile string) (deploymentConfig *v1alpha1.DeploymentConfig, err error)
}

type DeploymentConfig struct {
	DeploymentConfigAggregate
	deploymentConfigClient *mio.DeploymentConfig
	deployment             *kube.Deployment
	pipelineBuilder        builder.PipelineBuilder
	deploymentAggregate    DeploymentAggregate
}

func init() {
	app.Register(NewDeploymentConfigService)
}

func NewDeploymentConfigService(deploymentConfigClient *mio.DeploymentConfig, deployment *kube.Deployment, pipelineBuilder builder.PipelineBuilder, deploymentAggregate DeploymentAggregate) DeploymentConfigAggregate {
	return &DeploymentConfig{
		deploymentConfigClient: deploymentConfigClient,
		deployment:             deployment,
		pipelineBuilder:        pipelineBuilder,
		deploymentAggregate:    deploymentAggregate,
	}
}

func (d *DeploymentConfig) Template(cmd *command.DeploymentConfig) (deploymentConfig *v1alpha1.DeploymentConfig, err error) {
	log.Debug("build config templates create :%v", cmd)
	deploymentConfig = new(v1alpha1.DeploymentConfig)
	copier.Copy(deploymentConfig, cmd)
	deploymentConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.DeploymentConfigKind,
		APIVersion: constant.DeploymentConfigApiVersion,
	}
	deploymentConfig.ObjectMeta = v1.ObjectMeta{
		Name:      deploymentConfig.Name,
		Namespace: constant.TemplateDefaultNamespace,
	}
	deploymentConfig.Status.LastVersion = constant.InitLastVersion
	deployment, err := d.deploymentConfigClient.Get(deploymentConfig.Name, constant.TemplateDefaultNamespace)
	if err != nil {
		deploymentConfig, err = d.deploymentConfigClient.Create(deploymentConfig)
	} else {
		deployment.Spec = cmd.Spec
		deploymentConfig, err = d.deploymentConfigClient.Update(deploymentConfig.Name, constant.TemplateDefaultNamespace, deployment)
	}
	return
}

func (d *DeploymentConfig) Create(name, pipelineName, namespace, sourceType, version, buildVersion, profile string) (deploymentConfig *v1alpha1.DeploymentConfig, err error) {
	log.Debug("build config create name :%s, namespace : %s , sourceType : %s", name, namespace, sourceType)
	deploymentConfig = new(v1alpha1.DeploymentConfig)
	template, err := d.deploymentConfigClient.Get(sourceType, constant.TemplateDefaultNamespace)
	if err != nil {
		return nil, err
	}
	copier.Copy(deploymentConfig, template)
	deploymentConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.DeploymentConfigKind,
		APIVersion: constant.DeploymentConfigApiVersion,
	}
	deploymentConfig.ObjectMeta = v1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
	deploy, err := d.deploymentConfigClient.Get(name, namespace)
	if err == nil {
		deploy.Spec = template.Spec
		deploy.Spec.Profile = profile
		deploy.Status.LastVersion = deploy.Status.LastVersion + 1
		deploymentConfig, err = d.deploymentConfigClient.Update(name, namespace, deploy)
		log.Info("update deployment configs deploy :%v", deploymentConfig)
		log.Info("update deployment configs err :%v", err)
	} else {
		deploymentConfig.Status.LastVersion = constant.InitLastVersion
		deploymentConfig.Spec.Profile = profile
		deploymentConfig, err = d.deploymentConfigClient.Create(deploymentConfig)
	}
	d.deploymentAggregate.Create(deploymentConfig, pipelineName, version, buildVersion)
	return
}
