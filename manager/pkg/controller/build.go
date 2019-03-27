package controller

import (
	"github.com/jinzhu/copier"
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BuildController struct {
	at.RestController
	buildAggregate aggregate.BuildAggregate
}

func init() {
	app.Register(newBuildController)
}

func newBuildController(buildAggregate aggregate.BuildAggregate) *BuildController {
	return &BuildController{
		buildAggregate: buildAggregate,
	}
}

func (b *BuildController) Post(buildCommand *command.BuildConfig) (model.Response, error) {
	buildConfig := v1alpha1.BuildConfig{}
	copier.Copy(&buildConfig, buildCommand)
	build, err := b.buildAggregate.Create(&buildConfig, "", "v1")
	base := new(model.BaseResponse)
	base.SetData(build)
	return base, err
}

func (b *BuildController) GetByNameNamespace(name, namespace string) (model.Response, error) {
	build := &v1alpha1.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	err := b.buildAggregate.DeleteNode(build)
	base := new(model.BaseResponse)
	return base, err
}

func (b *BuildController) Delete() (model.Response, error) {
	build := &v1alpha1.Build{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"name": "hello-world",
			},
			Namespace: "demo",
		},
	}
	err := b.buildAggregate.DeleteNode(build)
	base := new(model.BaseResponse)
	return base, err
}
