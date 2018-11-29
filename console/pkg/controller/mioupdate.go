package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/mio/console/pkg/aggregate"
	"hidevops.io/mio/console/pkg/command"
)

type MioUpdateController struct {
	at.RestController
	mioUpdateAggregate aggregate.MioUpdateAggregate
}

func init() {
	app.Register(newMioUpdateController)
}

func newMioUpdateController(mioUpdateAggregate aggregate.MioUpdateAggregate) *MioUpdateController {
	return &MioUpdateController{
		mioUpdateAggregate: mioUpdateAggregate,
	}
}

func (c *MioUpdateController) Post(update *command.MioUpdate) (model.Response, error) {
	response := new(model.BaseResponse)
	name := update.Type + update.Arch
	err := c.mioUpdateAggregate.Add(name, update)
	return response, err
}

func (c *MioUpdateController) DeleteByTypeArch(types, arch string) (model.Response, error) {
	response := new(model.BaseResponse)
	name := types + arch
	err := c.mioUpdateAggregate.Delete(name)
	return response, err
}

func (c *MioUpdateController) GetByTypeArchVersion(types, arch, version string) (model.Response, error) {
	response := new(model.BaseResponse)
	name := types + arch
	update := new(command.MioUpdate)
	update, err := c.mioUpdateAggregate.Get(name)
	if version != update.Version {
		update.Enable = true
	}
	response.SetData(update)
	return response, err
}
