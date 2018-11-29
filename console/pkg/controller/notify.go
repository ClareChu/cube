package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/utils/copier"
	"hidevops.io/mio/console/pkg/aggregate"
	"hidevops.io/mio/console/pkg/command"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
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
