package builder

import (
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/console/pkg/command"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BuildNode interface {
	Start(node *command.DeployNode) (string, error)
	CreateServiceNode(node *command.ServiceNode) error
	DeleteDeployment(name, namespace string) error
	Update(name, namespace string) error
}

type BuildNodeImpl struct {
	BuildNode
	deployment *kube.Deployment
	service    *kube.Service
	replicaSet *kube.ReplicaSet
}

func init() {
	app.Register(newDeploymentConfig)
}

func newDeploymentConfig(deployment *kube.Deployment, service *kube.Service, replicaSet *kube.ReplicaSet) BuildNode {
	return &BuildNodeImpl{
		deployment: deployment,
		service:    service,
		replicaSet: replicaSet,
	}
}

func (s *BuildNodeImpl) Start(node *command.DeployNode) (string, error) {
	log.Infof("remote deploy: %v", node)
	d, err := s.deployment.DeployNode(&node.DeployData)
	return d, err
}

func (s *BuildNodeImpl) CreateServiceNode(node *command.ServiceNode) error {

	var ports []orch.Ports
	for _, port := range node.Ports {
		ports = append(ports, orch.Ports{
			Name: fmt.Sprintf("%d-tcp", port),
			Port: int32(port),
		})
	}

	err := s.service.Create(node.Name, node.NameSpace, ports)
	return err
}

func (s *BuildNodeImpl) DeleteDeployment(name, namespace string) error {
	deploy, err := s.deployment.Get(name, namespace, metav1.GetOptions{})
	replicas := int32(0)
	deploy.Spec.Replicas = &replicas
	err = s.deployment.Update(deploy)
	//TODO delete deployment
	Second := int64(0)
	options := &metav1.DeleteOptions{
		GracePeriodSeconds: &Second,
	}
	err = s.deployment.Delete(name, namespace, options)

	//TODO delete replica set
	option := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", name),
	}
	list, err := s.replicaSet.List(name, namespace, option)
	if  err != nil {
		return err
	}
	deleteOption := &metav1.DeleteOptions{
		GracePeriodSeconds: &Second,
	}
	err = s.replicaSet.Delete(list.Items[0].Name, namespace, deleteOption)
	//TODO delete service
	return err
}

func (s *BuildNodeImpl) Update(name, namespace string) error {
	deploy := &v1beta1.Deployment{}
	err := s.deployment.Update(deploy)
	return err
}
