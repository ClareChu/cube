package mio

import (
	"fmt"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	miov1alpha1 "hidevops.io/mio/pkg/client/clientset/versioned/typed/mio/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Notify struct {
	clientSet miov1alpha1.MioV1alpha1Interface
}

func NewNotify(clientSet miov1alpha1.MioV1alpha1Interface) *Notify {
	return &Notify{
		clientSet: clientSet,
	}
}

func (b *Notify) Create(notify *v1alpha1.Notify) (config *v1alpha1.Notify, err error) {
	log.Debugf("notify create : %v", notify.Name)
	config, err = b.clientSet.Notifies(notify.Namespace).Create(notify)
	if err != nil {
		return nil, err
	}
	return
}

func (b *Notify) Get(name, namespace string) (config *v1alpha1.Notify, err error) {
	log.Info(fmt.Sprintf("notify get app %s in namespace %s", name, namespace))
	result, err := b.clientSet.Notifies(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Notify) Delete(name, namespace string) error {
	log.Info(fmt.Sprintf("delete notify app %s in namespace %s", name, namespace))
	err := b.clientSet.Notifies(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (b *Notify) Update(name, namespace string, config *v1alpha1.Notify) (*v1alpha1.Notify, error) {
	log.Info(fmt.Sprintf("update notify app %s in namespace %s", name, namespace))
	result, err := b.clientSet.Notifies(namespace).Update(config)
	return result, err
}

func (b *Notify) List(namespace string, option v1.ListOptions) (*v1alpha1.NotifyList, error) {
	log.Info(fmt.Sprintf("list notify in namespace %s", namespace))
	result, err := b.clientSet.Notifies(namespace).List(option)
	return result, err
}

func (b *Notify) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s Notify，in the namespace %s", listOptions.LabelSelector, namespace))
	listOptions.Watch = true

	w, err := b.clientSet.Notifies(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
