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
	Create(name, pipelineName, namespace, sourceType, version, branch string) (buildConfig *v1alpha1.BuildConfig, err error)
	Delete(name, namespace string) error
}

type BuildConfig struct {
	BuildConfigAggregate
	buildConfigClient *cube.BuildConfig
	buildAggregate    BuildAggregate
}

func init() {
	app.Register(NewBuildConfigService)
}

func NewBuildConfigService(buildConfigClient *cube.BuildConfig, buildAggregate BuildAggregate) BuildConfigAggregate {
	return &BuildConfig{
		buildConfigClient: buildConfigClient,
		buildAggregate:    buildAggregate,
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

func (s *BuildConfig) Create(name, pipelineName, namespace, sourceType, version, branch string) (buildConfig *v1alpha1.BuildConfig, err error) {
	log.Debugf("build config create name :%v, namespace :%v", name, namespace)
	template, err := s.buildConfigClient.Get(sourceType, constant.TemplateDefaultNamespace)
	if err != nil {
		log.Errorf("get build config template err: %v", err)
		return nil, err
	}
	buildConfig, err = s.buildConfigClient.Get(name, namespace)
	buildConfigTemplate := new(v1alpha1.BuildConfig)
	copier.Copy(buildConfigTemplate, template)
	buildConfigTemplate.ObjectMeta = v1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels: map[string]string{
			constant.CodeType: sourceType,
		},
	}

	buildConfigTemplate.TypeMeta = v1.TypeMeta{
		Kind:       constant.BuildConfigKind,
		APIVersion: constant.BuildConfigApiVersion,
	}
	buildConfigTemplate.Spec.App = name
	buildConfigTemplate.Spec.CloneConfig.Branch = branch
	buildConfigTemplate.Spec.Tags = []string{template.Spec.DockerRegistry + "/" + namespace + "/" + name}
	//TODO 如果存在创建 buildconfig 不存在新建buildconfig 创建完buildconfig 新建
	if err != nil {
		buildConfigTemplate.Status.LastVersion = constant.InitLastVersion
		buildConfig, err = s.buildConfigClient.Create(buildConfigTemplate)
	} else {
		buildConfigTemplate.ObjectMeta = buildConfig.ObjectMeta
		buildConfigTemplate.Status.LastVersion = buildConfig.Status.LastVersion + 1
		buildConfig, err = s.buildConfigClient.Update(name, namespace, buildConfigTemplate)
	}
	if err != nil {
		log.Errorf("create build config :%v", err)
		return
	}
	//TODO 创建 build
	_, err = s.buildAggregate.Create(buildConfig, pipelineName, version)
	return
}
