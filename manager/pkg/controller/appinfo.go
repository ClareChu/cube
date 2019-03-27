package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hioak/starter/kube"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppInfoController struct {
	at.RestController
	pod *kube.Pod
}

func init() {
	app.Register(newAppInfoController)
}

func newAppInfoController(pod *kube.Pod) *AppInfoController {
	return &AppInfoController{
		pod: pod,
	}
}

func (a *AppInfoController) GetByNamespaceLabel(namespace, label string) (model.Response, error) {
	response := new(model.BaseResponse)
	if label == "all" {
		label = ""
	}
	podList, err := a.pod.GetPodList(namespace, meta_v1.ListOptions{
		LabelSelector: label,
	})

	response.SetData(podList)
	return response, err
}
