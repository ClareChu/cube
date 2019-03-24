package cube

import (
	"fmt"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	cubev1alpha1 "hidevops.io/cube/pkg/client/clientset/versioned/typed/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Build struct {
	clientSet cubev1alpha1.CubeV1alpha1Interface
}

func NewBuild(clientSet cubev1alpha1.CubeV1alpha1Interface) *Build {
	return &Build{
		clientSet: clientSet,
	}
}

func (b *Build) Create(build *v1alpha1.Build) (config *v1alpha1.Build, err error) {
	log.Debugf("build create : %v", build.Name)
	config, err = b.clientSet.Builds(build.Namespace).Create(build)
	if err != nil {
		return nil, err
	}
	return
}

func (b *Build) Get(name, namespace string) (config *v1alpha1.Build, err error) {
	log.Info(fmt.Sprintf("build get app %s in namespace %s", name, namespace))
	result, err := b.clientSet.Builds(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Build) Delete(name, namespace string) error {
	log.Info(fmt.Sprintf("delete build app %s in namespace %s", name, namespace))
	err := b.clientSet.Builds(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (b *Build) Update(name, namespace string, config *v1alpha1.Build) (*v1alpha1.Build, error) {
	log.Info(fmt.Sprintf("update build app %s in namespace %s", name, namespace))
	result, err := b.clientSet.Builds(namespace).Update(config)
	return result, err
}

func (b *Build) List(namespace string, option v1.ListOptions) (*v1alpha1.BuildList, error) {
	log.Info(fmt.Sprintf("list build in namespace %s", namespace))
	result, err := b.clientSet.Builds(namespace).List(option)
	return result, err
}

func (b *Build) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s Build，in the namespace %s", listOptions.LabelSelector, namespace))

	listOptions.Watch = true

	w, err := b.clientSet.Builds(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
