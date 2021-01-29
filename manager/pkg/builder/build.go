package builder

import (
	"fmt"
	"github.com/jinzhu/copier"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/service/client"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type BuildNode interface {
	Start(node *command.DeployNode) (string, error)
	CreateServiceNode(node *command.ServiceNode) error
	DeleteDeployment(name, namespace string) error
	Update(name, namespace string) error
	Delete(name, namespace string) (err error)
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
	clientSet, err := client.GetDefaultK8sClientSet()
	if err != nil {
		return err
	}
	err = Create(node.Name, node.App, node.NameSpace, ports, clientSet)
	return err
}

func Create(name, app, namespace string, ports interface{}, clientSet kubernetes.Interface) error {

	p := make([]corev1.ServicePort, 0)
	copier.Copy(&p, ports)

	// create service
	serviceSpec := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app":  name,
				"name": app,
			},
		},
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceTypeNodePort,
			Ports: p,
			Selector: map[string]string{
				"app": name,
			},
		},
	}

	svc, err := clientSet.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	switch {
	case err == nil:
		serviceSpec.ObjectMeta.ResourceVersion = svc.ObjectMeta.ResourceVersion
		serviceSpec.Spec.ClusterIP = svc.Spec.ClusterIP
		_, err = clientSet.CoreV1().Services(namespace).Update(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to update service: %s", err)
		}
		log.Info("service updated")
	case errors.IsNotFound(err):
		_, err = clientSet.CoreV1().Services(namespace).Create(serviceSpec)
		if err != nil {
			return fmt.Errorf("failed to create service")
		}
		log.Info("service created")
	default:
		return fmt.Errorf("upexected error: %s", err)
	}
	return nil
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

func (s *BuildNodeImpl) Delete(name, namespace string) (err error) {
	option := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", name),
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
