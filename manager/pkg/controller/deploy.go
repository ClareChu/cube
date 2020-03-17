package controller

import (
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type ReplicasController struct {
	at.RestController
	deployService service.DeployService
}

func init() {
	app.Register(newReplicasController)
}

func newReplicasController(deployService service.DeployService) *ReplicasController {
	return &ReplicasController{
		deployService: deployService,
	}
}


func (r *ReplicasController) Put(replicasRequest *service.ReplicasRequest) (response model.Response, err error) {
	response = new(model.BaseResponse)
	err = r.deployService.Update(replicasRequest)
	return response, err
}
