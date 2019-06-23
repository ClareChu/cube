package aggregate

import (
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/copier"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentConfigAggregate interface {
	Template(cmd *command.DeploymentConfig) (deploymentConfig *v1alpha1.DeploymentConfig, err error)
	Create(param *command.PipelineReqParams, buildVersion string) (deploymentConfig *v1alpha1.DeploymentConfig, err error)
	InitDeployConfig(deploy *v1alpha1.DeploymentConfig, template *v1alpha1.DeploymentConfig, param *command.PipelineReqParams)
}

type DeploymentConfig struct {
	DeploymentConfigAggregate
	deploymentConfigClient *cube.DeploymentConfig
	deployment             *kube.Deployment
	pipelineBuilder        builder.PipelineBuilder
	deploymentAggregate    DeploymentAggregate
}

func init() {
	app.Register(NewDeploymentConfigService)
}

func NewDeploymentConfigService(deploymentConfigClient *cube.DeploymentConfig, deployment *kube.Deployment, pipelineBuilder builder.PipelineBuilder, deploymentAggregate DeploymentAggregate) DeploymentConfigAggregate {
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

func (d *DeploymentConfig) Create(param *command.PipelineReqParams, buildVersion string) (deploymentConfig *v1alpha1.DeploymentConfig, err error) {
	log.Debugf("build config create name :%s, namespace : %s , sourceType : %s", param.Name, param.Namespace, param.EventType)
	deploymentConfig = new(v1alpha1.DeploymentConfig)
	template, err := d.deploymentConfigClient.Get(param.EventType, constant.TemplateDefaultNamespace)
	if err != nil {
		return nil, err
	}
	copier.Copy(deploymentConfig, template)
	deploymentConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.DeploymentConfigKind,
		APIVersion: constant.DeploymentConfigApiVersion,
	}
	deploymentConfig.ObjectMeta = v1.ObjectMeta{
		Name:      param.Name,
		Namespace: param.Namespace,
	}
	deploy, err := d.deploymentConfigClient.Get(param.Name, param.Namespace)
	if err == nil {
		d.InitDeployConfig(deploy, template, param)
		deploymentConfig, err = d.deploymentConfigClient.Update(param.Name, param.Namespace, deploy)
		log.Infof("update deployment configs deploy :%v", deploymentConfig)
	} else {
		deploymentConfig.Status.LastVersion = constant.InitLastVersion
		deploymentConfig.Spec.Profile = param.Profile
		deploymentConfig, err = d.deploymentConfigClient.Create(deploymentConfig)
	}
	d.deploymentAggregate.Create(deploymentConfig, param.PipelineName, param.Version, buildVersion)
	return
}

func (d *DeploymentConfig) InitDeployConfig(deploy *v1alpha1.DeploymentConfig, template *v1alpha1.DeploymentConfig, param *command.PipelineReqParams) {
	deploy.Spec = template.Spec
	deploy.Spec.Profile = param.Profile
	copier.Copy(&deploy.Spec.Container, &param.Container, copier.IgnoreEmptyValue)
	deploy.Spec.Container.Name = param.Name
	log.Info("---------env---------")
	log.Infof("env:%v", param.Container.Env)
	for _, e := range param.Container.Env {
		deploy.Spec.Container.Env = append(deploy.Spec.Container.Env, e)
	}
	deploy.Spec.Container.Env = append(append(deploy.Spec.Container.Env, corev1.EnvVar{Name: constant.AppName, Value: param.Name}), corev1.EnvVar{Name: constant.AppVersion, Value: param.Version})
	deploy.Status.LastVersion = deploy.Status.LastVersion + 1
}
