package crd

import (
	"hidevops.io/cube/operator/client"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineConfig struct {
	clientSet apiextension.Interface
}

func NewPipelineConfig(clientSet apiextension.Interface) CustomResourceDefinitions {
	return &PipelineConfig{
		clientSet: clientSet,
	}
}

func (ap *PipelineConfig) create() {
	aps := &v1beta1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: "pipelineconfigs.cube.io",
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   CubeGroup,
			Version: CubeVersion,
			Names: v1beta1.CustomResourceDefinitionNames{
				Kind:   "PipelineConfig",
				Plural: "pipelineconfigs",
			},
			Scope: v1beta1.NamespaceScoped,
		},
	}
	client.CreateCrd(ap.clientSet, aps)
}
