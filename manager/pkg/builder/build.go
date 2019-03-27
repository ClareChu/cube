package builder

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter"
	"hidevops.io/hioak/starter/kube"
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

	err := s.service.Create(node.Name, node.App, node.NameSpace, ports)
	return err
}

func (s *BuildNodeImpl) DeleteDeployment(name, namespace string) (err error) {
	option := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("name=%s", name),
	}
	deploys, err := s.deployment.List(namespace, option)
	opt := &metav1.DeleteOptions{}
	for _, deploy := range deploys.Items {

		err = s.deployment.Delete(deploy.Name, namespace, opt)
	}
	//TODO delete replica set

	list, err := s.replicaSet.List(name, namespace, option)
	if err != nil {
		return
	}
	for _, rs := range list.Items {
		err = s.replicaSet.Delete(rs.Name, namespace, opt)
		if err != nil {
			return
		}
	}
	return

}

func (s *BuildNodeImpl) Update(name, namespace string) error {
	deploy := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	err := s.deployment.Update(deploy)
	return err
}
