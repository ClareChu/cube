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

type PipelineConfigAggregate interface {
	NewPipelineConfigTemplate(PipelineConfigTemplate *command.PipelineConfigTemplate) (PipelineConfig *v1alpha1.PipelineConfig, err error)
	StartPipelineConfig(cmd *command.PipelineStart) (*v1alpha1.PipelineConfig, error)
	Get(name, namespace string) (*v1alpha1.PipelineConfig, error)
	Create(name, namespace string, pipelineConfig *v1alpha1.PipelineConfig) (*v1alpha1.PipelineConfig, error)
}

type PipelineConfig struct {
	PipelineConfigAggregate
	pipelineConfigClient *cube.PipelineConfig
	pipelineAggregate    PipelineAggregate
}

func init() {
	app.Register(NewPipelineConfigService)
}

func NewPipelineConfigService(pipelineConfigClient *cube.PipelineConfig, pipelineAggregate PipelineAggregate) PipelineConfigAggregate {
	return &PipelineConfig{
		pipelineConfigClient: pipelineConfigClient,
		pipelineAggregate:    pipelineAggregate,
	}
}

//新建 PipelineConfig 模版
func (p *PipelineConfig) NewPipelineConfigTemplate(pipelineConfigTemplate *command.PipelineConfigTemplate) (pipelineConfig *v1alpha1.PipelineConfig, err error) {
	log.Debug("build config create :%v", pipelineConfigTemplate)
	pipelineConfig = new(v1alpha1.PipelineConfig)
	copier.Copy(pipelineConfig, pipelineConfigTemplate)
	pipelineConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.PipelineConfigKind,
		APIVersion: constant.PipelineConfigApiVersion,
	}
	pipelineConfig.ObjectMeta = v1.ObjectMeta{
		Name:      pipelineConfig.Name,
		Namespace: constant.TemplateDefaultNamespace,
	}
	pipelineConfig.Status.LastVersion = 1
	pipeline, err := p.pipelineConfigClient.Get(pipelineConfig.Name, constant.TemplateDefaultNamespace)
	if err != nil {
		pipelineConfig, err = p.pipelineConfigClient.Create(pipelineConfig)
	} else {
		pipeline.Spec = pipelineConfigTemplate.Spec
		pipelineConfig, err = p.pipelineConfigClient.Update(pipelineConfig.Name, constant.TemplateDefaultNamespace, pipeline)
	}
	return
}

func (p *PipelineConfig) Get(name, namespace string) (*v1alpha1.PipelineConfig, error) {
	pipelineConfig, err := p.pipelineConfigClient.Get(name, namespace)
	return pipelineConfig, err
}

func (p *PipelineConfig) StartPipelineConfig(cmd *command.PipelineStart) (pipelineConfig *v1alpha1.PipelineConfig, err error) {
	log.Debugf("PipelineConfig get name: %v, namespace: %v", cmd.Name, cmd.Namespace)
	lastVersion := 1
	//TODO get pipeline template
	pipelineConfigTemplate, err := p.pipelineConfigClient.Get(cmd.SourceCode, constant.TemplateDefaultNamespace)
	if err != nil {
		log.Errorf("PipelineConfig get template : %v", err)
		return
	}
	pipelineConfig, err = p.pipelineConfigClient.Get(cmd.Name, cmd.Namespace)
	if err != nil {
		log.Errorf("PipelineConfig get err : %v", err)
		pipelineConfig = new(v1alpha1.PipelineConfig)
		copier.Copy(pipelineConfig, pipelineConfigTemplate)
		pipelineConfig.Status.LastVersion = lastVersion
		replaceProfile(cmd, pipelineConfig)
		pipelineConfig, err = p.Create(cmd.Name, cmd.Namespace, pipelineConfig)
	} else {
		lastVersion = pipelineConfig.Status.LastVersion + 1
		meta := pipelineConfig.ObjectMeta
		copier.Copy(pipelineConfig, pipelineConfigTemplate)
		pipelineConfig.ObjectMeta = meta
		replaceProfile(cmd, pipelineConfig)
		pipelineConfig.Status.LastVersion = lastVersion
		pipelineConfig, err = p.pipelineConfigClient.Update(cmd.Name, cmd.Namespace, pipelineConfig)
	}
	if err == nil {
		//TODO 	创建 pipeline
		_, err = p.pipelineAggregate.Create(pipelineConfig, cmd.SourceCode)
		return
	}
	return
}

func replaceProfile(cmd *command.PipelineStart, pipelineConfig *v1alpha1.PipelineConfig) {
	if cmd.Version != "" && cmd.Profile != "" && cmd.Branch != "" {
		pipelineConfig.Spec.Version = cmd.Version
		pipelineConfig.Spec.Profile = cmd.Profile
		pipelineConfig.Spec.Branch = cmd.Branch
	} else if cmd.Version != "" {
		pipelineConfig.Spec.Version = cmd.Version
	} else if cmd.Profile != "" {
		pipelineConfig.Spec.Profile = cmd.Profile
	} else if cmd.Branch != "" {
		pipelineConfig.Spec.Branch = cmd.Branch
	}
}

func (p *PipelineConfig) Create(name, namespace string, pipelineConfig *v1alpha1.PipelineConfig) (*v1alpha1.PipelineConfig, error) {
	pipelineConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.PipelineConfigKind,
		APIVersion: constant.PipelineConfigApiVersion,
	}
	pipelineConfig.ObjectMeta = v1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
	pipelineConfig.Spec.Namespace = namespace
	pipelineConfig.Spec.App = name
	return p.pipelineConfigClient.Create(pipelineConfig)
}
