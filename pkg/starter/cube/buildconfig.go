package cube

import (
	"fmt"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	cubev1alpha1 "hidevops.io/cube/pkg/client/clientset/versioned/typed/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type BuildConfig struct {
	clientSet cubev1alpha1.CubeV1alpha1Interface
}

func NewBuildConfig(clientSet cubev1alpha1.CubeV1alpha1Interface) *BuildConfig {
	return &BuildConfig{
		clientSet: clientSet,
	}
}

func (b *BuildConfig) Create(build *v1alpha1.BuildConfig) (config *v1alpha1.BuildConfig, err error) {
	log.Debugf("buildConfig create : %v", build.Name)
	config, err = b.clientSet.BuildConfigs(build.Namespace).Create(build)
	if err != nil {
		return nil, err
	}
	return
}

func (b *BuildConfig) Get(name, namespace string) (config *v1alpha1.BuildConfig, err error) {
	log.Infof("get buildConfig app %s in namespace %s", name, namespace)
	result, err := b.clientSet.BuildConfigs(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *BuildConfig) Delete(name, namespace string) error {
	log.Info(fmt.Sprintf("delete buildConfig app %s in namespace %s", name, namespace))
	err := b.clientSet.BuildConfigs(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (b *BuildConfig) Update(name, namespace string, config *v1alpha1.BuildConfig) (*v1alpha1.BuildConfig, error) {
	log.Info(fmt.Sprintf("update buildConfig app %s in namespace %s", name, namespace))
	result, err := b.clientSet.BuildConfigs(namespace).Update(config)
	return result, err
}

func (b *BuildConfig) List(namespace string, option v1.ListOptions) (*v1alpha1.BuildConfigList, error) {
	log.Info(fmt.Sprintf("list buildConfig in namespace %s", namespace))
	result, err := b.clientSet.BuildConfigs(namespace).List(option)
	return result, err
}

func (b *BuildConfig) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s BuildConfig，in the namespace %s", listOptions.LabelSelector, namespace))
	listOptions.Watch = true

	w, err := b.clientSet.BuildConfigs(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
