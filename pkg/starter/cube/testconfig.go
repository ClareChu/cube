package cube

import (
	"fmt"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	cubev1alpha1 "hidevops.io/cube/pkg/client/clientset/versioned/typed/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type TestConfig struct {
	clientSet cubev1alpha1.CubeV1alpha1Interface
}

func NewTestConfig(clientSet cubev1alpha1.CubeV1alpha1Interface) *TestConfig {
	return &TestConfig{
		clientSet: clientSet,
	}
}

func (d *TestConfig) Create(testConfig *v1alpha1.TestConfig) (result *v1alpha1.TestConfig, err error) {
	log.Debugf("deployConfig create : %v", testConfig.Name)
	result, err = d.clientSet.TestConfigs(testConfig.Namespace).Create(testConfig)
	if err != nil {
		return nil, err
	}
	return
}

func (d *TestConfig) Get(name, namespace string) (config *v1alpha1.TestConfig, err error) {
	log.Info(fmt.Sprintf("get testConfig %s in namespace %s", name, namespace))
	result, err := d.clientSet.TestConfigs(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *TestConfig) Delete(name, namespace string) error {
	log.Info(fmt.Sprintf("delete testConfig %s in namespace %s", name, namespace))
	err := d.clientSet.TestConfigs(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (d *TestConfig) Update(name, namespace string, testConfig *v1alpha1.TestConfig) (*v1alpha1.TestConfig, error) {
	log.Info(fmt.Sprintf("update testConfig %s in namespace %s", name, namespace))
	result, err := d.clientSet.TestConfigs(namespace).Update(testConfig)
	return result, err
}

func (d *TestConfig) List(namespace string, option v1.ListOptions) (*v1alpha1.TestConfigList, error) {
	log.Info(fmt.Sprintf("list deployConfig in namespace %s", namespace))
	result, err := d.clientSet.TestConfigs(namespace).List(option)
	return result, err
}

func (d *TestConfig) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s testconfig，in the namespace %s", listOptions.LabelSelector, namespace))
	w, err := d.clientSet.TestConfigs(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
