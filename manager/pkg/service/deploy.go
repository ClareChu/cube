package service

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hioak/starter/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeployService interface {
	Update(replicasRequest *ReplicasRequest) (err error)
}

type DeployServiceImpl struct {
	DeployService
	deployment *kube.Deployment
}

type ReplicasRequest struct {
	model.RequestBody
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas  *int32 `json:"replicas"`
}

func init() {
	app.Register(newDeployCommand)
}

func newDeployCommand(deployment *kube.Deployment) DeployService {
	return &DeployServiceImpl{
		deployment: deployment,
	}
}

func (a *DeployServiceImpl) Update(replicasRequest *ReplicasRequest) (err error) {
	option := metav1.GetOptions{}
	res, err := a.deployment.Get(replicasRequest.Name, replicasRequest.Namespace, option)
	if err != nil {
		log.Errorf("get deployment error:%v", err)
		return
	}

	res.Spec.Replicas = replicasRequest.Replicas
	return a.deployment.Update(res)
}
