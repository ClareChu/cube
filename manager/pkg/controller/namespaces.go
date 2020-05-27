package controller

import (
	"hidevops.io/cube/manager/pkg/service/client"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceController struct {
	at.RestController
}

func init() {
	app.Register(newNamespaceController)
}

func newNamespaceController() *NamespaceController {
	return &NamespaceController{
	}
}

func (c *NamespaceController) DeleteByName(namespace string) (model.Response, error) {
	response := new(model.BaseResponse)
	clientSet, err := client.GetDefaultK8sClientSet()
	if err != nil {
		return response, err
	}
	ops := &metav1.DeleteOptions{}
	err = clientSet.CoreV1().Namespaces().Delete(namespace, ops)
	return response, err
}
