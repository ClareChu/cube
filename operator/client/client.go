package client

import (
	"hidevops.io/hiboot/pkg/log"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

func CreateCrd(clientSet apiextension.Interface, crd *v1beta1.CustomResourceDefinition) {
	_, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil {
		log.Errorf("create CustomResourceDefinition error :%v", err)
		return
	}
	log.Infof("create CustomResourceDefinition :%v success", crd.Name)

}
