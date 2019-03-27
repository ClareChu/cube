package aggregate

import (
	"fmt"
	"github.com/kevholditch/gokong"
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
	Create(name, pipelineName, namespace, sourceType, version, profile string) (gatewayConfig *v1alpha1.GatewayConfig, err error)
	Gateway(gatewayConfig *v1alpha1.GatewayConfig) (err error)
}

type GatewayConfig struct {
	GatewayConfigAggregate
	gatewayConfigClient *cube.GatewayConfig
	pipelineBuilder     builder.PipelineBuilder
}

func init() {
	app.Register(NewGatewayService)
}

func NewGatewayService(gatewayConfigClient *cube.GatewayConfig, pipelineBuilder builder.PipelineBuilder) GatewayConfigAggregate {
	return &GatewayConfig{
		gatewayConfigClient: gatewayConfigClient,
		pipelineBuilder:     pipelineBuilder,
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

func (s *GatewayConfig) Create(name, pipelineName, namespace, sourceType, version, profile string) (gatewayConfig *v1alpha1.GatewayConfig, err error) {
	log.Debugf("gateway config create name :%s, namespace : %s , sourceType : %s", name, namespace, sourceType)
	phase := constant.Success
	project := namespace
	gatewayConfig = new(v1alpha1.GatewayConfig)
	if profile != "" {
		namespace = fmt.Sprintf("%s-%s", namespace, profile)
	}
	template, err := s.gatewayConfigClient.Get(sourceType, constant.TemplateDefaultNamespace)
	if err != nil {
		return nil, err
	}
	template.Name = fmt.Sprintf("%s-%s", namespace, name)
	template.Spec.UpstreamUrl = fmt.Sprintf("http://%s.%s.svc:8080", name, namespace)
	uri := fmt.Sprintf("/%s/%s", namespace, name)
	uri = strings.Replace(uri, "-", "/", -1)
	template.Spec.Uris = []string{uri}
	copier.Copy(gatewayConfig, template)
	gatewayConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.DeploymentConfigKind,
		APIVersion: constant.DeploymentConfigApiVersion,
	}
	gatewayConfig.ObjectMeta = v1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels: map[string]string{
			constant.PipelineConfigName: pipelineName,
		},
	}

	gateway, err := s.gatewayConfigClient.Get(name, namespace)
	if err == nil {
		gateway.Spec = template.Spec
		gatewayConfig, err = s.gatewayConfigClient.Update(name, namespace, gateway)
	} else {
		gatewayConfig, err = s.gatewayConfigClient.Create(gatewayConfig)
	}
	err = s.Gateway(gatewayConfig)
	if err != nil {
		log.Errorf("create gateway err : %v", err)
		phase = constant.Fail
	}
	err = s.pipelineBuilder.Update(pipelineName, project, constant.CreateService, phase, "")
	return
}

func (s *GatewayConfig) Gateway(gatewayConfig *v1alpha1.GatewayConfig) (err error) {
	apiRequest := &gokong.ApiRequest{
		Name:                   gatewayConfig.Name,
		Hosts:                  gatewayConfig.Spec.Hosts,
		Uris:                   gatewayConfig.Spec.Uris,
		UpstreamUrl:            gatewayConfig.Spec.UpstreamUrl,
		StripUri:               gatewayConfig.Spec.StripUri,
		PreserveHost:           gatewayConfig.Spec.PreserveHost,
		Retries:                gatewayConfig.Spec.Retries,
		UpstreamConnectTimeout: gatewayConfig.Spec.UpstreamConnectTimeout,
		UpstreamSendTimeout:    gatewayConfig.Spec.UpstreamSendTimeout,
		UpstreamReadTimeout:    gatewayConfig.Spec.UpstreamReadTimeout,
		HttpsOnly:              gatewayConfig.Spec.HttpsOnly,
		HttpIfTerminated:       gatewayConfig.Spec.HttpIfTerminated,
	}
	config := &gokong.Config{
		HostAddress: gatewayConfig.Spec.KongAdminUrl,
	}
	_, err = gokong.NewClient(config).Apis().Create(apiRequest)
	return
}
