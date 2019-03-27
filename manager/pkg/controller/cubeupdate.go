package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type CubeUpdateController struct {
	at.RestController
	cubeUpdateAggregate aggregate.CubeUpdateAggregate
}

func init() {
	app.Register(newCubeUpdateController)
}

func newCubeUpdateController(cubeUpdateAggregate aggregate.CubeUpdateAggregate) *CubeUpdateController {
	return &CubeUpdateController{
		cubeUpdateAggregate: cubeUpdateAggregate,
	}
}

func (c *CubeUpdateController) Post(update *command.CubeUpdate) (model.Response, error) {
	response := new(model.BaseResponse)
	name := update.Type + update.Arch
	err := c.cubeUpdateAggregate.Add(name, update)
	return response, err
}

func (c *CubeUpdateController) DeleteByTypeArch(types, arch string) (model.Response, error) {
	response := new(model.BaseResponse)
	name := types + arch
	err := c.cubeUpdateAggregate.Delete(name)
	return response, err
}

func (c *CubeUpdateController) GetByTypeArchVersion(types, arch, version string) (model.Response, error) {
	response := new(model.BaseResponse)
	name := types + arch
	update := new(command.CubeUpdate)
	update, err := c.cubeUpdateAggregate.Get(name)
	if version != update.Version {
		update.Enable = true
	}
	response.SetData(update)
	return response, err
}
