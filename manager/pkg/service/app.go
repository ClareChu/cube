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
	"reflect"
)

type AppService interface {
	Create(cmd *command.PipelineStart) (app *v1alpha1.App, err error)
	Get(name, namespace string) (app *v1alpha1.App, err error)
	Delete(name, namespace string) (err error)
	Update(name, namespace string, cmd *command.PipelineStart) (app *v1alpha1.App, err error)
	Init(cmd *command.PipelineStart) (similarity bool, err error)
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
	if cmd.Name == "" {
		cmd.Name = cmd.AppRoot
	}
	if cmd.Namespace == "" {
		cmd.Namespace = cmd.Project
	}
	name := fmt.Sprintf("%s-%s-%s", cmd.Project, cmd.AppRoot, cmd.Version)
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
	app, err = a.appClient.Get(name, namespace)
	return
}

func (a *AppServiceImpl) Delete(name, namespace string) (err error) {
	log.Infof("delete app name: %s, namespace %s", name, namespace)
	err = a.appClient.Delete(name, namespace)
	return
}

func (a *AppServiceImpl) Compare(name string, oldApp, newApp v1alpha1.AppSpec) (similarity bool, err error) {
	//多次登录导致token不一致的问题
	oldApp.Token = ""
	newApp.Token = ""
	if !reflect.DeepEqual(oldApp, newApp) {
		log.Debugf("*** oldApp: %v", oldApp)
		log.Debugf("*** newApp: %v", newApp)
		log.Infof("update app")
		var cmd command.PipelineStart
		err = copier.Copy(&cmd, newApp)
		if err != nil {
			return false, err
		}
		_, err = a.Create(&cmd)
		return false, nil
	}
	log.Debugf("*** app similarity")
	return true, nil
}

//Init cmd value validate
func (a *AppServiceImpl) Init(cmd *command.PipelineStart) (similarity bool, err error) {
	//todo 通过URL部署项目
	name := fmt.Sprintf("%s-%s-%s", cmd.Project, cmd.AppRoot, cmd.Version)
	app, err := a.Get(name, constant.TemplateDefaultNamespace)
	if err == nil {
		//TODO 查看多次入参是否相似
		var spec v1alpha1.AppSpec
		err = copier.Copy(&spec, app.Spec)
		err = copier.Copy(app.Spec, cmd, copier.IgnoreEmptyValue)
		similarity, err := a.Compare(name, spec, app.Spec)
		err = copier.Copy(cmd, app.Spec)

		return similarity, err
	} else {
		// TODO create app
		app, err = a.Create(cmd)
		log.Errorf("create app err :%v", err)
	}
	return false, nil
}
