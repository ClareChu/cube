package aggregate

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/aggregate/dispatch"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/manager/utils"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/copier"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

type GatewayConfigAggregate interface {
	Template(cmd *command.GatewayConfig) (gatewayConfig *v1alpha1.GatewayConfig, err error)
	Create(params *command.PipelineReqParams) (err error)
}

type GatewayConfig struct {
	GatewayConfigAggregate
	gatewayConfigClient      *cube.GatewayConfig
	pipelineBuilder          builder.PipelineBuilder
	gatewayAggregate         GatewayAggregate
	pipelineFactoryInterface dispatch.PipelineFactoryInterface
}

const GatewayEvent = "gateway"

func init() {
	app.Register(NewGatewayService)
}

func NewGatewayService(gatewayConfigClient *cube.GatewayConfig,
	pipelineBuilder builder.PipelineBuilder,
	gatewayAggregate GatewayAggregate,
	pipelineFactoryInterface dispatch.PipelineFactoryInterface) GatewayConfigAggregate {
	gatewayConfig := &GatewayConfig{
		gatewayConfigClient:      gatewayConfigClient,
		pipelineBuilder:          pipelineBuilder,
		gatewayAggregate:         gatewayAggregate,
		pipelineFactoryInterface: pipelineFactoryInterface,
	}
	pipelineFactoryInterface.Add(GatewayEvent, gatewayConfig)
	return gatewayConfig
}

func (s *GatewayConfig) Template(cmd *command.GatewayConfig) (gatewayConfig *v1alpha1.GatewayConfig, err error) {
	log.Debug("build config templates create :%v", cmd)
	gatewayConfig = new(v1alpha1.GatewayConfig)
	err = copier.Copy(gatewayConfig, cmd)
	if err != nil {
		log.Errorf("copy is template error: %v", err)
		return
	}
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

func (s *GatewayConfig) Create(params *command.PipelineReqParams) (err error) {
	log.Debugf("gateway config create name :%s, namespace : %s , sourceType : %s", params.Name, params.Namespace, params.EventType)
	phase := constant.Success
	template, err := s.gatewayConfigClient.Get(params.EventType, constant.TemplateDefaultNamespace)
	if err != nil {
		log.Errorf("get gateway config template err :%v", err)
		return err
	}
	annotations := make(map[string]string)
	for k, v := range template.Annotations {
		log.Debugf("Annotations debug v:%v", v)
		reg := regexp.MustCompile(`\$\{(.*?)\}`)
		match := reg.MatchString(v)
		log.Debugf("match debug v:%v", match)
		if match {
			value := utils.Regx(v, "${", "}", params)
			log.Debugf("value debug v:%v", value)
			annotations[k] = value
		} else {
			annotations[k] = v
		}
	}

	if len(params.Ingress) != 0 {
		for _, ing := range params.Ingress {
			uri := ""
			container := strings.Index(ing.Path, "/")
			if container != 0 {
				uri = fmt.Sprintf("/%s", ing.Path)
			} else {
				uri = ing.Path
			}
			uri = fmt.Sprintf("%s%s", uri, template.Spec.RegexPath)


			gatewayConfig := &v1alpha1.GatewayConfig{
				TypeMeta: v1.TypeMeta{
					Kind:       constant.GatewayConfigKind,
					APIVersion: constant.GatewayConfigApiVersion,
				},
				ObjectMeta: v1.ObjectMeta{
					Name:        ing.Name,
					Namespace:   params.Namespace,
					Annotations: template.Annotations,
					Labels: map[string]string{
						constant.PipelineConfigName: ing.Name,
						constant.Namespace:          params.Namespace,
					},
				},
				Spec: v1alpha1.GatewaySpec{
					Hosts: []string{ing.Domain},
					Port:  ing.Port,
					Uris:  []string{uri},
				},
			}
			err = s.gatewayAggregate.Create(gatewayConfig)
			if err != nil {
				log.Errorf("create gateway err : %v", err)
				phase = constant.Fail
				break
			}
		}
		err = s.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.Gateway, phase, "")
		return
	}

	template.Name = fmt.Sprintf("%s-%s", params.Namespace, params.Name)
	uri := fmt.Sprintf("/%s/%s%s", params.Namespace, params.Name, template.Spec.RegexPath)
	template.Spec.Uris = []string{uri}
	gatewayConfig := &v1alpha1.GatewayConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       constant.GatewayConfigKind,
			APIVersion: constant.GatewayConfigApiVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:        params.Name,
			Namespace:   params.Namespace,
			Annotations: annotations,
			Labels: map[string]string{
				constant.PipelineConfigName: params.PipelineName,
				constant.Namespace:          params.Namespace,
			},
		},
		Spec: template.Spec,
	}
	gateway, err := s.gatewayConfigClient.Get(params.Name, params.Namespace)
	if err == nil {
		gateway.Spec = template.Spec
		gateway.ObjectMeta.Annotations = annotations
		gatewayConfig, err = s.gatewayConfigClient.Update(params.Name, params.Namespace, gateway)
	} else {
		gatewayConfig, err = s.gatewayConfigClient.Create(gatewayConfig)
	}
	err = s.gatewayAggregate.Create(gatewayConfig)
	if err != nil {
		log.Errorf("create gateway err : %v", err)
		phase = constant.Fail
	}
	err = s.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.Gateway, phase, "")
	return
}
