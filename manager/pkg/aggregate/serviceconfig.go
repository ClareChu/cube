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
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceConfigAggregate interface {
	Template(cmd *command.ServiceConfig) (serviceConfig *v1alpha1.ServiceConfig, err error)
	Create(params *command.PipelineReqParams) (serviceConfig *v1alpha1.ServiceConfig, err error)
	CreateServices(name, namespace string, sc []v1alpha1.Service) (err error)
	DeleteService(name, namespace string) (err error)
	Delete(name, namespace string) (err error)
}

type ServiceConfig struct {
	ServiceConfigAggregate
	serviceConfigClient *cube.ServiceConfig
	service             *kube.Service
	pipelineBuilder     builder.PipelineBuilder
}

func init() {
	app.Register(NewServiceConfigService)
}

func NewServiceConfigService(serviceConfigClient *cube.ServiceConfig, service *kube.Service, pipelineBuilder builder.PipelineBuilder) ServiceConfigAggregate {
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

func (s *ServiceConfig) Create(params *command.PipelineReqParams) (serviceConfig *v1alpha1.ServiceConfig, err error) {
	log.Debugf("create service name :%s, namespace : %s , sourceType : %s", params.Name, params.Namespace, params.EventType)
	phase := constant.Success
	serviceConfig = new(v1alpha1.ServiceConfig)
	template, err := s.serviceConfigClient.Get(params.EventType, constant.TemplateDefaultNamespace)
	if err != nil {
		log.Infof("create service err : %v", err)
		return nil, err
	}
	copier.Copy(serviceConfig, template)
	serviceConfig.TypeMeta = v1.TypeMeta{
		Kind:       constant.DeploymentConfigKind,
		APIVersion: constant.DeploymentConfigApiVersion,
	}
	serviceConfig.ObjectMeta = v1.ObjectMeta{
		Name:      params.Name,
		Namespace: params.Namespace,
		Labels: map[string]string{
			constant.CodeType: params.EventType,
		},
	}
	deploy, err := s.serviceConfigClient.Get(params.Name, params.Namespace)
	if err == nil {
		deploy.Spec = template.Spec
		serviceConfig, err = s.serviceConfigClient.Update(params.Name, params.Namespace, deploy)
	} else {
		serviceConfig, err = s.serviceConfigClient.Create(serviceConfig)
	}
	if len(params.Services) == 0 {
		err = s.CreateService(serviceConfig)
	} else {
		err = s.CreateServices(params.Name, params.Namespace, params.Services)
	}

	if err != nil {
		phase = constant.Fail
		log.Errorf("create service name %v err : %v", params.Name, err)
	}
	log.Info("create service success")
	err = s.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.CreateService, phase, "")
	return
}

func (s *ServiceConfig) CreateService(serviceConfig *v1alpha1.ServiceConfig) (err error) {
	err = s.service.CreateService(serviceConfig.Name, serviceConfig.Namespace, serviceConfig.Spec.Ports)
	return
}

func (s *ServiceConfig) CreateServices(name, namespace string, sc []v1alpha1.Service) (err error) {
	for _, svc := range sc {
		sv, err := s.service.Get(svc.Name, namespace)
		serviceSpec := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: svc.Name,

				Labels: map[string]string{
					"app": name,
				},
			},
			Spec: corev1.ServiceSpec{
				Type:  corev1.ServiceTypeClusterIP,
				Ports: svc.Ports,
				Selector: map[string]string{
					"app": name,
				},
			},
		}
		if err != nil {
			_, err = s.service.CreateSvc(namespace, serviceSpec)
			if err != nil {
				return err
			}
		} else {
			//TODO update svc
			serviceSpec.ObjectMeta.ResourceVersion = sv.ObjectMeta.ResourceVersion
			serviceSpec.Spec.ClusterIP = sv.Spec.ClusterIP
			_, err = s.service.Update(namespace, serviceSpec)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ServiceConfig) DeleteService(name, namespace string) (err error) {
	opt := v1.ListOptions{
		LabelSelector: fmt.Sprintf("name=%s", name),
	}
	list, err := s.service.List(namespace, opt)
	for _, service := range list.Items {
		err = s.service.Delete(service.Name, namespace)
		if err != nil {
			return
		}
	}
	return
}

func (s *ServiceConfig) Delete(name, namespace string) (err error) {
	err = s.service.Delete(name, namespace)
	return
}
