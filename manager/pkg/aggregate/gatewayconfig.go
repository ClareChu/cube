package aggregate

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/copier"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type GatewayConfigAggregate interface {
	Template(cmd *command.GatewayConfig) (gatewayConfig *v1alpha1.GatewayConfig, err error)
	Create(params *command.PipelineReqParams) (gatewayConfig *v1alpha1.GatewayConfig, err error)
}

type GatewayConfig struct {
	GatewayConfigAggregate
	gatewayConfigClient *cube.GatewayConfig
	pipelineBuilder     builder.PipelineBuilder
	gatewayAggregate    GatewayAggregate
}

func init() {
	app.Register(NewGatewayService)
}

func NewGatewayService(gatewayConfigClient *cube.GatewayConfig, pipelineBuilder builder.PipelineBuilder, gatewayAggregate GatewayAggregate) GatewayConfigAggregate {
	return &GatewayConfig{
		gatewayConfigClient: gatewayConfigClient,
		pipelineBuilder:     pipelineBuilder,
		gatewayAggregate:    gatewayAggregate,
	}
}

func (s *GatewayConfig) Template(cmd *command.GatewayConfig) (gatewayConfig *v1alpha1.GatewayConfig, err error) {
	log.Debug("build config templates create :%v", cmd)
	gatewayConfig = new(v1alpha1.GatewayConfig)
	copier.Copy(gatewayConfig, cmd)
	gatewayConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.GatewayConfigKind,
		APIVersion: constant.GatewayConfigApiVersion,
	}
	gatewayConfig.ObjectMeta = v1.ObjectMeta{
		Name:      gatewayConfig.Name,
		Namespace: constant.TemplateDefaultNamespace,
	}
	service, err := s.gatewayConfigClient.Get(gatewayConfig.Name, constant.TemplateDefaultNamespace)
	if err != nil {
		gatewayConfig, err = s.gatewayConfigClient.Create(gatewayConfig)
	} else {
		service.Spec = cmd.Spec
		gatewayConfig, err = s.gatewayConfigClient.Update(gatewayConfig.Name, constant.TemplateDefaultNamespace, service)
	}
	return
}

func (s *GatewayConfig) Create(params *command.PipelineReqParams) (gatewayConfig *v1alpha1.GatewayConfig, err error) {
	log.Debugf("gateway config create name :%s, namespace : %s , sourceType : %s", params.Name, params.Namespace, params.EventType)
	phase := constant.Success
	project := params.Namespace
	gatewayConfig = new(v1alpha1.GatewayConfig)
	template, err := s.gatewayConfigClient.Get(params.EventType, constant.TemplateDefaultNamespace)
	if err != nil {
		return nil, err
	}
	template.Name = fmt.Sprintf("%s-%s", params.Namespace, params.Name)
	template.Spec.UpstreamUrl = fmt.Sprintf("http://%s.%s.svc:8080", params.Name, params.Namespace)
	uri := fmt.Sprintf("/%s/%s", project, params.Name)
	uri = strings.Replace(uri, "-", "/", -1)
	template.Spec.Uris = []string{uri}
	copier.Copy(gatewayConfig, template)
	gatewayConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.GatewayConfigKind,
		APIVersion: constant.GatewayConfigApiVersion,
	}
	gatewayConfig.ObjectMeta = v1.ObjectMeta{
		Name:      params.Name,
		Namespace: params.Namespace,
		Labels: map[string]string{
			constant.PipelineConfigName: params.PipelineName,
			constant.Namespace:          project,
		},
	}
	gateway, err := s.gatewayConfigClient.Get(params.Name, params.Namespace)
	if err == nil {
		gateway.Spec = template.Spec
		gatewayConfig, err = s.gatewayConfigClient.Update(params.Name, params.Namespace, gateway)
	} else {
		gatewayConfig, err = s.gatewayConfigClient.Create(gatewayConfig)
	}
	err = s.gatewayAggregate.Create(gatewayConfig)
	if err != nil {
		log.Errorf("create gateway err : %v", err)
		phase = constant.Fail
	}

	err = s.pipelineBuilder.Update(params.PipelineName, project, constant.CreateService, phase, "")
	return
}
