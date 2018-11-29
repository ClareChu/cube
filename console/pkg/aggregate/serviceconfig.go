package aggregate

import (
	"fmt"
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

type ServiceConfigAggregate interface {
	Template(cmd *command.ServiceConfig) (serviceConfig *v1alpha1.ServiceConfig, err error)
	Create(name, pipelineName, namespace, sourceType, version, profile string) (serviceConfig *v1alpha1.ServiceConfig, err error)
	DeleteService(name, namespace string) (err error)
}

type ServiceConfig struct {
	ServiceConfigAggregate
	serviceConfigClient *mio.ServiceConfig
	service             *kube.Service
	pipelineBuilder     builder.PipelineBuilder
}

func init() {
	app.Register(NewServiceConfigService)
}

func NewServiceConfigService(serviceConfigClient *mio.ServiceConfig, service *kube.Service, pipelineBuilder builder.PipelineBuilder) ServiceConfigAggregate {
	return &ServiceConfig{
		serviceConfigClient: serviceConfigClient,
		service:             service,
		pipelineBuilder:     pipelineBuilder,
	}
}

func (s *ServiceConfig) Template(cmd *command.ServiceConfig) (serviceConfig *v1alpha1.ServiceConfig, err error) {
	log.Debug("build config templates create :%v", cmd)
	serviceConfig = new(v1alpha1.ServiceConfig)
	copier.Copy(serviceConfig, cmd)
	serviceConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.ServiceConfigKind,
		APIVersion: constant.ServiceConfigApiVersion,
	}
	serviceConfig.ObjectMeta = v1.ObjectMeta{
		Name:      serviceConfig.Name,
		Namespace: constant.TemplateDefaultNamespace,
		Labels:    cmd.ObjectMeta.Labels,
	}
	service, err := s.serviceConfigClient.Get(serviceConfig.Name, constant.TemplateDefaultNamespace)
	if err != nil {
		serviceConfig, err = s.serviceConfigClient.Create(serviceConfig)
	} else {
		service.Spec = cmd.Spec
		serviceConfig, err = s.serviceConfigClient.Update(serviceConfig.Name, constant.TemplateDefaultNamespace, service)
	}
	return
}

func (s *ServiceConfig) Create(name, pipelineName, namespace, sourceType, version, profile string) (serviceConfig *v1alpha1.ServiceConfig, err error) {
	log.Debug("build config create name :%s, namespace : %s , sourceType : %s", name, namespace, sourceType)
	phase := constant.Success
	serviceConfig = new(v1alpha1.ServiceConfig)
	template, err := s.serviceConfigClient.Get(sourceType, constant.TemplateDefaultNamespace)
	if err != nil {
		log.Infof("create service err : %v", err)
		return nil, err
	}
	if profile != "" {
		namespace = fmt.Sprintf("%s-%s", namespace, profile)
	}
	copier.Copy(serviceConfig, template)
	serviceConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.DeploymentConfigKind,
		APIVersion: constant.DeploymentConfigApiVersion,
	}
	serviceConfig.ObjectMeta = v1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels: map[string]string{
			constant.CodeType: sourceType,
		},
	}
	deploy, err := s.serviceConfigClient.Get(name, namespace)
	if err == nil {
		deploy.Spec = template.Spec
		serviceConfig, err = s.serviceConfigClient.Update(name, namespace, deploy)
	} else {
		serviceConfig, err = s.serviceConfigClient.Create(serviceConfig)
	}
	err = s.CreateService(serviceConfig)
	if err != nil {
		phase = constant.Fail
		log.Errorf("create service name %v err : %v", name, err)
	}
	err = s.pipelineBuilder.Update(pipelineName, namespace, constant.CreateService, phase, "")
	return
}

func (s *ServiceConfig) CreateService(serviceConfig *v1alpha1.ServiceConfig) (err error) {
	err = s.service.CreateService(serviceConfig.Name, serviceConfig.Namespace, serviceConfig.Spec.Ports)
	return
}

func (s *ServiceConfig) DeleteService(name, namespace string) (err error) {
	err = s.service.Delete(name, namespace)
	return err
}
