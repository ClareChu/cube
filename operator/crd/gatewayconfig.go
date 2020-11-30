package crd

import (
	"hidevops.io/cube/operator/client"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GatewayConfig struct {
	clientSet apiextension.Interface
}

func NewGatewayConfig(clientSet apiextension.Interface) CustomResourceDefinitions {
	return &GatewayConfig{
		clientSet: clientSet,
	}
}

func (ap *GatewayConfig) create() {
	aps := &v1beta1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: "gatewayconfigs.cube.io",
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   CubeGroup,
			Version: CubeVersion,
			Names: v1beta1.CustomResourceDefinitionNames{
				Kind:   "GatewayConfig",
				Plural: "gatewayconfigs",
			},
			Scope: v1beta1.NamespaceScoped,
		},
	}
	client.CreateCrd(ap.clientSet, aps)
}
