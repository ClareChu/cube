package mio

import (
	"fmt"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	miov1alpha1 "hidevops.io/mio/pkg/client/clientset/versioned/typed/mio/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Tests struct {
	clientSet miov1alpha1.MioV1alpha1Interface
}

func NewTestses(clientSet miov1alpha1.MioV1alpha1Interface) *Tests {
	return &Tests{
		clientSet: clientSet,
	}
}

func (d *Tests) Create(tests *v1alpha1.Tests) (result *v1alpha1.Tests, err error) {
	log.Debugf("deploy create : %v", tests.Name)
	result, err = d.clientSet.Testses(tests.Namespace).Create(tests)
	if err != nil {
		return nil, err
	}
	return
}

func (d *Tests) Get(name, namespace string) (result *v1alpha1.Tests, err error) {
	log.Info(fmt.Sprintf("get tests %s in namespace %s", name, namespace))
	result, err = d.clientSet.Testses(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *Tests) Delete(name, namespace string) error {
	log.Info(fmt.Sprintf("delete tests %s in namespace %s", name, namespace))
	err := d.clientSet.Testses(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (d *Tests) Update(name, namespace string, tests *v1alpha1.Tests) (*v1alpha1.Tests, error) {
	log.Info(fmt.Sprintf("update tests %s in namespace %s", name, namespace))
	result, err := d.clientSet.Testses(namespace).Update(tests)
	return result, err
}

func (d *Tests) List(namespace string, option v1.ListOptions) (*v1alpha1.TestsList, error) {
	log.Info(fmt.Sprintf("list tests in namespace %s", namespace))
	result, err := d.clientSet.Testses(namespace).List(option)
	fmt.Println("Error", err)
	return result, err
}

func (d *Tests) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s Tests，in the namespace %s", listOptions.LabelSelector, namespace))
	w, err := d.clientSet.Testses(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
