package mio

import (
	"fmt"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	miov1alpha1 "hidevops.io/mio/pkg/client/clientset/versioned/typed/mio/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type DeploymentConfig struct {
	clientSet miov1alpha1.MioV1alpha1Interface
}

func NewDeploymentConfig(clientSet miov1alpha1.MioV1alpha1Interface) *DeploymentConfig {
	return &DeploymentConfig{
		clientSet: clientSet,
	}
}

func (d *DeploymentConfig) Create(deploymentConfig *v1alpha1.DeploymentConfig) (result *v1alpha1.DeploymentConfig, err error) {
	log.Debugf("deployConfig create : %v", deploymentConfig.Name)
	result, err = d.clientSet.DeploymentConfigs(deploymentConfig.Namespace).Create(deploymentConfig)
	if err != nil {
		return nil, err
	}
	return
}

func (d *DeploymentConfig) Get(name, namespace string) (config *v1alpha1.DeploymentConfig, err error) {
	log.Info(fmt.Sprintf("get deployConfig app %s in namespace %s", name, namespace))
	result, err := d.clientSet.DeploymentConfigs(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DeploymentConfig) Delete(name, namespace string) error {
	log.Info(fmt.Sprintf("delete deployConfig app %s in namespace %s", name, namespace))
	err := d.clientSet.DeploymentConfigs(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (d *DeploymentConfig) Update(name, namespace string, deploymentConfig *v1alpha1.DeploymentConfig) (*v1alpha1.DeploymentConfig, error) {
	log.Info(fmt.Sprintf("update deployConfig app %s in namespace %s", name, namespace))
	result, err := d.clientSet.DeploymentConfigs(namespace).Update(deploymentConfig)
	return result, err
}

func (d *DeploymentConfig) List(namespace string, option v1.ListOptions) (*v1alpha1.DeploymentConfigList, error) {
	log.Info(fmt.Sprintf("list deployConfig in namespace %s", namespace))
	result, err := d.clientSet.DeploymentConfigs(namespace).List(option)
	return result, err
}

func (d *DeploymentConfig) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s DeploymentConfig，in the namespace %s", listOptions.LabelSelector, namespace))

	listOptions.Watch = true

	w, err := d.clientSet.DeploymentConfigs(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
