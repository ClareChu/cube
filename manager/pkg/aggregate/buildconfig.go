package aggregate

import (
	"github.com/jinzhu/copier"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BuildConfigAggregate interface {
	Template(buildConfigTemplate *command.BuildConfig) (buildConfig *v1alpha1.BuildConfig, err error)
	Create(params *command.PipelineReqParams) (buildConfig *v1alpha1.BuildConfig, err error)
	Delete(name, namespace string) error
	InitConfig(buildConfigTemplate *v1alpha1.BuildConfig, params *command.PipelineReqParams, template *v1alpha1.BuildConfig)
}

type BuildConfig struct {
	BuildConfigAggregate
	buildConfigClient   *cube.BuildConfig
	buildAggregate      BuildAggregate
	configMapsAggregate ConfigMapsAggregate
}

func init() {
	app.Register(NewBuildConfigService)
}

func NewBuildConfigService(buildConfigClient *cube.BuildConfig, buildAggregate BuildAggregate, configMapsAggregate ConfigMapsAggregate) BuildConfigAggregate {
	return &BuildConfig{
		buildConfigClient:   buildConfigClient,
		buildAggregate:      buildAggregate,
		configMapsAggregate: configMapsAggregate,
	}
}

func (s *BuildConfig) Delete(name, namespace string) error {
	err := s.buildConfigClient.Delete(name, namespace)
	log.Errorf("delete build config name:%v , namespace :%v , err: %v", name, namespace, err)
	return err
}

//新建 buildConfig 模版
func (s *BuildConfig) Template(buildConfigTemplate *command.BuildConfig) (buildConfig *v1alpha1.BuildConfig, err error) {
	log.Debug("build config create :%v", buildConfigTemplate)
	buildConfig = new(v1alpha1.BuildConfig)
	copier.Copy(buildConfig, buildConfigTemplate)
	buildConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.BuildConfigKind,
		APIVersion: constant.BuildApiVersion,
	}
	buildConfig.ObjectMeta = v1.ObjectMeta{
		Name:      buildConfig.Name,
		Namespace: constant.TemplateDefaultNamespace,
	}
	buildConfig.Status.LastVersion = constant.InitLastVersion
	build, err := s.buildConfigClient.Get(buildConfig.Name, constant.TemplateDefaultNamespace)
	if err != nil {
		buildConfig, err = s.buildConfigClient.Create(buildConfig)
	} else {
		build.Spec = buildConfigTemplate.Spec
		buildConfig, err = s.buildConfigClient.Update(buildConfig.Name, constant.TemplateDefaultNamespace, build)
	}
	return
}

func (s *BuildConfig) Create(params *command.PipelineReqParams) (buildConfig *v1alpha1.BuildConfig, err error) {
	log.Debugf("build config create name :%v, namespace :%v", params.Name, params.Namespace)
	template, err := s.buildConfigClient.Get(params.EventType, constant.TemplateDefaultNamespace)
	if err != nil {
		log.Errorf("get build config template err: %v", err)
		return nil, err
	}
	config, err := s.configMapsAggregate.Get(constant.DockerConstant, constant.TemplateDefaultNamespace)
	if err != nil {
		log.Errorf("get configMaps err :%v", err)
	}
	template.Spec.DockerAuthConfig.Username = config.Data[constant.Username]
	template.Spec.DockerAuthConfig.Password = config.Data[constant.Password]
	template.Spec.DockerRegistry = config.Data[constant.DockerRegistry]
	buildConfig, err = s.buildConfigClient.Get(params.Name, params.Namespace)
	buildConfigTemplate := new(v1alpha1.BuildConfig)
	copier.Copy(buildConfigTemplate, template)
	buildConfigTemplate.ObjectMeta = v1.ObjectMeta{
		Name:      params.Name,
		Namespace: params.Namespace,
		Labels: map[string]string{
			constant.CodeType:   params.EventType,
			constant.AppName:    params.Name,
			constant.AppVersion: params.Version,
		},
	}

	buildConfigTemplate.TypeMeta = v1.TypeMeta{
		Kind:       constant.BuildConfigKind,
		APIVersion: constant.BuildConfigApiVersion,
	}
	s.InitConfig(buildConfigTemplate, params, template)
	//TODO 如果存在创建 buildConfig 不存在新建 buildConfig 创建完 buildConfig 新建
	if err != nil {
		buildConfigTemplate.Status.LastVersion = constant.InitLastVersion
		buildConfig, err = s.buildConfigClient.Create(buildConfigTemplate)
	} else {
		buildConfigTemplate.ObjectMeta = buildConfig.ObjectMeta
		buildConfigTemplate.Status.LastVersion = buildConfig.Status.LastVersion + 1
		buildConfig, err = s.buildConfigClient.Update(params.Name, params.Namespace, buildConfigTemplate)
	}
	if err != nil {
		log.Errorf("create build config :%v", err)
		return
	}
	//TODO 创建 build
	_, err = s.buildAggregate.Create(buildConfig, params.PipelineName, params.Version)
	return
}

func (s *BuildConfig) InitConfig(buildConfigTemplate *v1alpha1.BuildConfig, params *command.PipelineReqParams, template *v1alpha1.BuildConfig) {
	buildConfigTemplate.Spec.App = params.Name
	buildConfigTemplate.Spec.CloneConfig.Branch = params.Branch
	buildConfigTemplate.Spec.Context = params.Context
	buildConfigTemplate.Spec.AppRoot = params.AppRoot
	buildConfigTemplate.Spec.Project = params.Project
	buildConfigTemplate.Spec.Tags = []string{template.Spec.DockerRegistry + "/" + params.Namespace + "/" + params.Name}
}
