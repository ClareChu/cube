package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/utils/copier"
)

type NotifyController struct {
	at.RestController
	notifyAggregate aggregate.NotifyAggregate
}

func init() {
	app.Register(newNotifyController)
}

func newNotifyController(notifyAggregate aggregate.NotifyAggregate) *NotifyController {
	return &NotifyController{
		notifyAggregate: notifyAggregate,
	}
}

func (c *NotifyController) Post(cmd *command.Notify) (model.Response, error) {
	notify := &v1alpha1.Notify{}
	copier.Copy(notify, cmd)
	c.notifyAggregate.Create(notify)
	return nil, nil
}

func (c *NotifyController) GetByNameNamespace(name, namespace string) (model.Response, error) {

	return nil, nil
}
