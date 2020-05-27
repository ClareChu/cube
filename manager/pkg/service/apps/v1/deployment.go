package v1

import (
	"hidevops.io/hiboot/pkg/log"
	appsV1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Deployment struct {
	ClientSet *kubernetes.Clientset
}

func (d *Deployment) Create(deploy *appsV1.Deployment) (*appsV1.Deployment, error) {
	log.Debugf("create deployment name :%v, namespace :%v", deploy.Name, deploy.Namespace)
	dpm, err := d.ClientSet.AppsV1().Deployments(deploy.Namespace).Create(deploy)
	return dpm, err
}

func (d *Deployment) Delete(name, namespace string, option *metav1.DeleteOptions) error {
	log.Debugf("delete deployment name :%v, namespace :%v", name, namespace)
	err := d.ClientSet.AppsV1().Deployments(namespace).Delete(name, option)
	return err
}

func (d *Deployment) Update(deployment *appsV1.Deployment) error {
	log.Debugf("update deployment name :%v, namespace :%v", deployment.Name, deployment.Namespace)
	_, err := d.ClientSet.AppsV1().Deployments(deployment.Namespace).Update(deployment)
	return err
}

func (d *Deployment) Get(name, namespace string, option metav1.GetOptions) (*appsV1.Deployment, error) {
	log.Debugf("get deployment name :%v, namespace :%v", name, namespace)
	deploy, err := d.ClientSet.AppsV1().Deployments(namespace).Get(name, option)
	return deploy, err
}

func (d *Deployment) List(namespace string, option metav1.ListOptions) (*appsV1.DeploymentList, error) {
	log.Debugf("get deployment namespace :%v", namespace)
	deploys, err := d.ClientSet.AppsV1().Deployments(namespace).List(option)
	return deploys, err
}
