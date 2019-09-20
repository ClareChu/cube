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
	corev1 "k8s.io/api/core/v1"
	extensionsV1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sort"
)

type DeploymentConfigAggregate interface {
	Template(cmd *command.DeploymentConfig) (deploymentConfig *v1alpha1.DeploymentConfig, err error)
	Create(param *command.PipelineReqParams, buildVersion string) (deploymentConfig *v1alpha1.DeploymentConfig, err error)
	InitDeployConfig(deploy *v1alpha1.DeploymentConfig, template *v1alpha1.DeploymentConfig, param *command.PipelineReqParams)
}

type DeploymentConfig struct {
	DeploymentConfigAggregate
	deploymentConfigClient *cube.DeploymentConfig
	pipelineBuilder        builder.PipelineBuilder
	deploymentAggregate    DeploymentAggregate
}

func init() {
	app.Register(NewDeploymentConfigService)
}

func NewDeploymentConfigService(deploymentConfigClient *cube.DeploymentConfig, pipelineBuilder builder.PipelineBuilder, deploymentAggregate DeploymentAggregate) DeploymentConfigAggregate {
	return &DeploymentConfig{
		deploymentConfigClient: deploymentConfigClient,
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
		d.InitDeployConfig(deploymentConfig, template, param)
		deploymentConfig, err = d.deploymentConfigClient.Create(deploymentConfig)
	}
	_, err = d.deploymentAggregate.Create(deploymentConfig, param.PipelineName, param.Version, buildVersion)
	return nil, err
}

//TODO 初始化 init deploy config
func (d *DeploymentConfig) InitDeployConfig(deploy *v1alpha1.DeploymentConfig, template *v1alpha1.DeploymentConfig, param *command.PipelineReqParams) {
	deploy.Spec = template.Spec
	deploy.Spec.Profile = param.Profile
	copier.Copy(&deploy.Spec.Container, &param.Container, copier.IgnoreEmptyValue)
	envs := map[string]string{}
	for _, e := range template.Spec.Container.Env {
		envs[e.Name] = e.Value
	}
	for _, e := range deploy.Spec.Container.Env {
		envs[e.Name] = e.Value
	}
	envs[constant.AppName] = param.Name
	envs[constant.AppVersion] = param.Version
	envVars := []corev1.EnvVar{}
	for k, v := range envs {
		envVar := corev1.EnvVar{
			Name:  k,
			Value: v,
		}
		envVars = append(envVars, envVar)
	}
	if param.Volumes.Name != "" {
		deploy.Spec.Strategy.Type = extensionsV1beta1.RecreateDeploymentStrategyType
		deploy.Spec.Volumes = []corev1.Volume{
			corev1.Volume{
				Name: param.Volumes.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: param.Volumes.Name,
					},
				},
			},
		}
		deploy.Spec.Container.VolumeMounts = []corev1.VolumeMount{
			corev1.VolumeMount{
				Name:      param.Volumes.Name,
				MountPath: param.Volumes.Workspace,
			},
		}
	} else {
		deploy.Spec.Strategy = extensionsV1beta1.DeploymentStrategy{
			Type: extensionsV1beta1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &extensionsV1beta1.RollingUpdateDeployment{
				MaxUnavailable: &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(0),
				},
				MaxSurge: &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(1),
				},
			},
		}
	}
	if param.InitContainer.Name != "" {
		deploy.Spec.InitContainer = corev1.Container{
			Name:            param.InitContainer.Name,
			Image:           param.InitContainer.Image,
			Command:         param.InitContainer.Command,
			ImagePullPolicy: corev1.PullAlways,
			Env:             param.InitContainer.Env,
			VolumeMounts: []corev1.VolumeMount{
				corev1.VolumeMount{
					Name:      param.Volumes.Name,
					MountPath: param.Volumes.Workspace,
				},
			},
		}
	}
	log.Debugf("*********** InitialDelaySeconds: %v", deploy.Spec.ReadinessProbe)
	if deploy.Spec.ReadinessProbe.InitialDelaySeconds != 0 {
		deploy.Spec.Container.ReadinessProbe = &deploy.Spec.ReadinessProbe
	}

	if deploy.Spec.LivenessProbe.InitialDelaySeconds != 0 {
		deploy.Spec.Container.LivenessProbe = &deploy.Spec.LivenessProbe
	}
	deploy.Spec.ForceUpdate = param.ForceUpdate
	deploy.Spec.Container.Name = param.Name
	deploy.Spec.Container.Env = sortEnv(envVars)
	deploy.Status.LastVersion = deploy.Status.LastVersion + 1
}

func sortEnv(envs []corev1.EnvVar) []corev1.EnvVar {
	var s []string
	var es []corev1.EnvVar
	for _, env := range envs {
		s = append(s, env.Name)
	}
	sort.Strings(s)
	for _, d := range s {
		for _, env := range envs {
			if d == env.Name {
				es = append(es, env)
			}
		}
	}
	return es
}
