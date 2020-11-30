package crd

import (
	"hidevops.io/cube/manager/pkg/service/client"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
)

const (
	CubeVersion = "v1alpha1"
	CubeGroup   = "cube.io"
)

type InitCustomResourceDefinitions func(clientSet apiextension.Interface) CustomResourceDefinitions

type CustomResourceDefinitions interface {
	create()
}

var definitionsInterface = []InitCustomResourceDefinitions{
	NewApplication,
	NewDeployment,
	NewDeploymentConfig,
	NewGatewayConfig,
	NewPipeline,
	NewPipelineConfig,
	NewServiceConfig,
	NewSourceConfig,
}

type CustomResourceDefinition struct {
	clientSet apiextension.Interface
}

func InitCRD() (crd *CustomResourceDefinition, err error) {
	clientSet, err := client.GetDefaultApiExtensionClientSet()
	if err != nil {
		return
	}
	return &CustomResourceDefinition{clientSet: clientSet}, nil
}

func fakeInitCRD() (crd *CustomResourceDefinition, err error) {
	clientSet := fake.NewSimpleClientset()
	return &CustomResourceDefinition{clientSet: clientSet}, nil
}

// 初始化crd的所有的资源
func (c *CustomResourceDefinition) Run() {
	for _, d := range definitionsInterface {
		d(c.clientSet).create()
	}
}
