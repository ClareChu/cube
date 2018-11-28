package mio

import (
	"fmt"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	miov1 "hidevops.io/mio/pkg/client/clientset/versioned/typed/mio/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type GatewayConfig struct {
	clientSet miov1.MioV1alpha1Interface
}

func NewGatewayConfig(clientSet miov1.MioV1alpha1Interface) *GatewayConfig {
	return &GatewayConfig{
		clientSet: clientSet,
	}
}

func (s *GatewayConfig) Create(gatewayConfigs *v1alpha1.GatewayConfig) (config *v1alpha1.GatewayConfig, err error) {
	log.Debugf("gatewayConfigs create : %v", gatewayConfigs.Name)
	config, err = s.clientSet.GatewayConfigs(gatewayConfigs.Namespace).Create(gatewayConfigs)
	if err != nil {
		return nil, err
	}
	return
}

func (s *GatewayConfig) Get(name, namespace string) (config *v1alpha1.GatewayConfig, err error) {
	log.Info("get GatewayConfigs ", name)
	result, err := s.clientSet.GatewayConfigs(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *GatewayConfig) Delete(name, namespace string) error {
	log.Info("delete GatewayConfig ", name)
	err := s.clientSet.GatewayConfigs(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (s *GatewayConfig) Update(name, namespace string, config *v1alpha1.GatewayConfig) (*v1alpha1.GatewayConfig, error) {
	log.Info("update GatewayConfig ", name)
	result, err := s.clientSet.GatewayConfigs(namespace).Update(config)
	return result, err
}

func (s *GatewayConfig) List(namespace string, option v1.ListOptions) (*v1alpha1.GatewayConfigList, error) {
	log.Info(fmt.Sprintf("list GatewayConfig in namespace %s", namespace))
	result, err := s.clientSet.GatewayConfigs(namespace).List(option)
	return result, err
}

func (s *GatewayConfig) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s GatewayConfig，in the namespace %s", listOptions.LabelSelector, namespace))

	listOptions.Watch = true

	w, err := s.clientSet.GatewayConfigs(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
