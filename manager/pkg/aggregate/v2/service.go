package v2

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/aggregate/dispatch"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceInterface interface {
	Handle(build *v1alpha1.Build) error
	DeleteService(name, namespace string) (err error)
}

type Service struct {
	ServiceInterface
	buildNode          builder.BuildNode
	build              BuildInterface
	service            *kube.Service
	buildPackInterface dispatch.BuildPackInterface
}

func init() {
	app.Register(NewServiceService)
}

const CreateService = "createService"

func NewServiceService(buildNode builder.BuildNode,
	build BuildInterface,
	service *kube.Service,
	buildPackInterface dispatch.BuildPackInterface) ServiceInterface {
	svc := &Service{
		buildNode:          buildNode,
		build:              build,
		service:            service,
		buildPackInterface: buildPackInterface,
	}
	buildPackInterface.Add(CreateService, svc)
	return svc
}

func (s *Service) Handle(build *v1alpha1.Build) error {
	log.Infof("aggregate create service %v", build)
	phase := constant.Success
	cmd := &command.ServiceNode{
		DeployData: kube.DeployData{
			Name:      build.Name,
			App:       build.ObjectMeta.Labels["name"],
			NameSpace: build.Namespace,
			Replicas:  build.Spec.DeployData.Replicas,
			Labels:    build.Spec.DeployData.Labels,
			Image:     build.Spec.BaseImage,
			Ports:     build.Spec.DeployData.Ports,
		},
	}
	err := s.buildNode.CreateServiceNode(cmd)
	if err != nil {
		log.Errorf("deploy hinode err :%v", err)
		phase = constant.Fail
		return err
	}
	err = s.build.Update(build, constant.CreateService, phase)
	return err
}

func (s *Service) DeleteService(name, namespace string) (err error) {
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
