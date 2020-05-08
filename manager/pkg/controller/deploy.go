package controller

import (
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
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

const DefaultVersion = "v1"
/*
func (r *ReplicasController) Put(replicasRequest *service.ReplicasRequest, ctx context.Context) (response model.Response, err error) {
	token := ctx.GetHeader("Authorization")
	response = new(model.BaseResponse)
	if token == "" {
		return response, errors.New("Permission denied authorization  token is not found")
	}
	replicasRequest.Token = token
	if replicasRequest.Version == "" {
		replicasRequest.Version = DefaultVersion
	}
	replicasRequest.App = replicasRequest.Name + "-" + replicasRequest.Version
	err = r.deployService.Put(replicasRequest)
	return response, err
}

*/
//Put 现在回调不需要加鉴权了
func (r *ReplicasController) Put(replicasRequest *service.ReplicasRequest, ctx context.Context) (response model.Response, err error) {
	response = new(model.BaseResponse)
	if replicasRequest.Version == "" {
		replicasRequest.Version = DefaultVersion
	}
	replicasRequest.App = replicasRequest.Name + "-" + replicasRequest.Version
	err = r.deployService.Put(replicasRequest)
	return response, err
}

func (r *ReplicasController) Update(replicasRequest *service.ReplicasRequest) (response model.Response, err error) {
	response = new(model.BaseResponse)
	if replicasRequest.Version == "" {
		replicasRequest.Version = DefaultVersion
	}
	replicasRequest.Name = replicasRequest.Name + "-" + replicasRequest.Version
	err = r.deployService.Update(replicasRequest)
	return response, err
}
