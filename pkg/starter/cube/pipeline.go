package cube

import (
	"fmt"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	cubev1 "hidevops.io/cube/pkg/client/clientset/versioned/typed/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Pipeline struct {
	clientSet cubev1.CubeV1alpha1Interface
}

func NewPipeline(clientSet cubev1.CubeV1alpha1Interface) *Pipeline {
	return &Pipeline{
		clientSet: clientSet,
	}
}

func (b *Pipeline) Create(pipeline *v1alpha1.Pipeline) (config *v1alpha1.Pipeline, err error) {
	log.Debugf("pipeline create : %v", pipeline.Name)
	config, err = b.clientSet.Pipelines(pipeline.Namespace).Create(pipeline)
	if err != nil {
		return nil, err
	}
	return
}

func (b *Pipeline) Get(name, namespace string) (config *v1alpha1.Pipeline, err error) {
	log.Info("get pipeline ", name)
	result, err := b.clientSet.Pipelines(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Pipeline) Delete(name, namespace string) error {
	log.Info("delete pipeline ", name)
	err := b.clientSet.Pipelines(namespace).Delete(name, &v1.DeleteOptions{})
	return err
}

func (b *Pipeline) Update(name, namespace string, config *v1alpha1.Pipeline) (*v1alpha1.Pipeline, error) {
	log.Info("update pipeline ", name)
	result, err := b.clientSet.Pipelines(namespace).Update(config)
	return result, err
}

func (b *Pipeline) List(namespace string, option v1.ListOptions) (*v1alpha1.PipelineList, error) {
	log.Info(fmt.Sprintf("list pipeline in namespace %s", namespace))
	result, err := b.clientSet.Pipelines(namespace).List(option)
	return result, err
}

func (b *Pipeline) Watch(listOptions v1.ListOptions, namespace string) (watch.Interface, error) {
	log.Info(fmt.Sprintf("watch label for %s Pipeline，in the namespace %s", listOptions.LabelSelector, namespace))

	listOptions.Watch = true

	w, err := b.clientSet.Pipelines(namespace).Watch(listOptions)
	if err != nil {
		return nil, err
	}
	return w, nil
}
