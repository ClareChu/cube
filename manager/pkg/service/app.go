package service

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/copier"
)

type AppService interface {
	Create(cmd *command.PipelineStart) (app *v1alpha1.App, err error)
	Get(name, namespace string) (app *v1alpha1.App, err error)
	Delete(name, namespace string) (err error)
	Update(name, namespace string, cmd *command.PipelineStart) (app *v1alpha1.App, err error)
}

type AppServiceImpl struct {
	AppService
	appClient *cube.App
}

func init() {
	app.Register(newAppCommand)
}

func newAppCommand(appClient *cube.App) AppService {
	return &AppServiceImpl{
		appClient: appClient,
	}
}

func (a *AppServiceImpl) Create(cmd *command.PipelineStart) (app *v1alpha1.App, err error) {
	name := fmt.Sprintf("%s-%s-%s", cmd.Project, cmd.Name, cmd.Version)
	namespace := constant.TemplateDefaultNamespace
	app1, err := a.Get(name, namespace)
	app = &v1alpha1.App{}
	copier.Copy(&app.Spec, cmd)
	app.Name = name
	app.Namespace = namespace
	if err != nil {
		app, err = a.appClient.Create(app)
		return
	}
	app1.Spec = app.Spec
	app, err = a.appClient.Update(name, namespace, app1)
	return
}

func (a *AppServiceImpl) Get(name, namespace string) (app *v1alpha1.App, err error) {
	log.Infof("get app name: %s, namespace %s", name, namespace)
	app, err = a.appClient.Get(name, namespace)
	return
}

func (a *AppServiceImpl) Delete(name, namespace string) (err error) {
	log.Infof("delete app name: %s, namespace %s", name, namespace)
	err = a.appClient.Delete(name, namespace)
	return
}
