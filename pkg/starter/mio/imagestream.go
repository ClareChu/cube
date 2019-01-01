package mio

import (
	"fmt"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	miov1 "hidevops.io/mio/pkg/client/clientset/versioned/typed/mio/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type ImageStream struct {
	clientSet miov1.MioV1alpha1Interface
}

func NewImageStream(clientSet miov1.MioV1alpha1Interface) *ImageStream {
	return &ImageStream{
		clientSet: clientSet,
	}
}

func (s *ImageStream) Create(imageStream *v1alpha1.ImageStream) (image *v1alpha1.ImageStream, err error) {
	log.Debugf("image stream create : %v", imageStream.Name)
	image, err = s.clientSet.ImageStreams(imageStream.Namespace).Create(imageStream)
	if err != nil {
		return nil, err
	}
	return
}

func (s *ImageStream) Get(name, namespace string) (image *v1alpha1.ImageStream, err error) {
	log.Info("get ImageStream ", name)
	result, err := s.clientSet.ImageStreams(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ImageStream) Delete(name, namespace string) error {
	log.Info("delete ImageStream ", name)
	err := s.clientSet.ImageStreams(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (s *ImageStream) Update(name, namespace string, imageStream *v1alpha1.ImageStream) (*v1alpha1.ImageStream, error) {
	log.Info("update ImageStreams ", name)
	result, err := s.clientSet.ImageStreams(namespace).Update(imageStream)
	return result, err
}

func (s *ImageStream) List(namespace string, option v1.ListOptions) (*v1alpha1.ImageStreamList, error) {
	log.Info(fmt.Sprintf("list ImageStream in namespace %s", namespace))
	result, err := s.clientSet.ImageStreams(namespace).List(option)
	return result, err
}

func (s *ImageStream) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s ImageStream，in the namespace %s", listOptions.LabelSelector, namespace))

	listOptions.Watch = true

	w, err := s.clientSet.ImageStreams(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
