package mio

import (
	"fmt"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	miov1alpha1 "hidevops.io/mio/pkg/client/clientset/versioned/typed/mio/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Deployment struct {
	clientSet miov1alpha1.MioV1alpha1Interface
}

func NewDeployment(clientSet miov1alpha1.MioV1alpha1Interface) *Deployment {
	return &Deployment{
		clientSet: clientSet,
	}
}

func (d *Deployment) Create(deployment *v1alpha1.Deployment) (result *v1alpha1.Deployment, err error) {
	log.Debugf("deploy create : %v", deployment.Name)
	result, err = d.clientSet.Deployments(deployment.Namespace).Create(deployment)
	if err != nil {
		return nil, err
	}
	return
}

func (d *Deployment) Get(name, namespace string) (result *v1alpha1.Deployment, err error) {
	log.Info(fmt.Sprintf("get deploy app %s in namespace %s", name, namespace))
	result, err = d.clientSet.Deployments(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *Deployment) Delete(name, namespace string) error {
	log.Info(fmt.Sprintf("delete deploy app %s in namespace %s", name, namespace))
	err := d.clientSet.Deployments(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (d *Deployment) Update(name, namespace string, deployment *v1alpha1.Deployment) (*v1alpha1.Deployment, error) {
	log.Info(fmt.Sprintf("update deploy app %s in namespace %s", name, namespace))
	result, err := d.clientSet.Deployments(namespace).Update(deployment)
	return result, err
}

func (d *Deployment) List(namespace string, option v1.ListOptions) (*v1alpha1.DeploymentList, error) {
	log.Info(fmt.Sprintf("list deploy in namespace %s", namespace))
	result, err := d.clientSet.Deployments(namespace).List(option)
	return result, err
}

func (d *Deployment) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s Deployment，in the namespace %s", listOptions.LabelSelector, namespace))
	listOptions.Watch = true

	w, err := d.clientSet.Deployments(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
