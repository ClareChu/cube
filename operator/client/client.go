package client

import (
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned"
	"hidevops.io/hiboot/pkg/log"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

func CreateCrd(clientSet apiextension.Interface, crd *v1beta1.CustomResourceDefinition) {
	_, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil {
		log.Errorf("create CustomResourceDefinition error :%v", err)
		return
	}
	log.Infof("create CustomResourceDefinition :%v success", crd.Name)
}

func CreateRestClient(client rest.Interface, obj interface{}, resource string) {
	var result runtime.Object
	err := client.Post().
		Resource(resource).
		Body(obj).
		Do().
		Into(result)
	if err != nil {
		log.Errorf("create  error :%v", err)
		return
	}
	log.Infof("create  :%v success")
}

func CreatDeploymentConfig(client versioned.Interface, deploy *v1alpha1.DeploymentConfig) {
	_, err := client.CubeV1alpha1().DeploymentConfigs(deploy.Namespace).Create(deploy)
	if err != nil {
		log.Errorf("creat DeploymentConfig error :%v", err)
		return
	}
	log.Infof("create DeploymentConfig :%v success", deploy.Name)
}

func CreatGatewayConfig(client versioned.Interface, deploy *v1alpha1.GatewayConfig) {
	_, err := client.CubeV1alpha1().GatewayConfigs(deploy.Namespace).Create(deploy)
	if err != nil {
		log.Errorf("creat GatewayConfigs error :%v", err)
		return
	}
	log.Infof("create GatewayConfigs :%v success", deploy.Name)
}

func CreatPipelineConfig(client versioned.Interface, pipeline *v1alpha1.PipelineConfig) {
	_, err := client.CubeV1alpha1().PipelineConfigs(pipeline.Namespace).Create(pipeline)
	if err != nil {
		log.Errorf("creat PipelineConfigs error :%v", err)
		return
	}
	log.Infof("create PipelineConfigs :%v success", pipeline.Name)
}

func CreatServiceConfig(client versioned.Interface, pipeline *v1alpha1.ServiceConfig) {
	_, err := client.CubeV1alpha1().ServiceConfigs(pipeline.Namespace).Create(pipeline)
	if err != nil {
		log.Errorf("creat ServiceConfigs error :%v", err)
		return
	}
	log.Infof("create ServiceConfigs :%v success", pipeline.Name)
}
