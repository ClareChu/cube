package controller

import (
	"hidevops.io/cube/manager/pkg/aggregate"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
)

type ConfigurationController struct {
	at.RestController
	configMapsAggregate aggregate.ConfigMapsAggregate
}

func init() {
	app.Register(newInitConfigurationController)
}

func newInitConfigurationController(configMapsAggregate aggregate.ConfigMapsAggregate) *ConfigurationController {
	return &ConfigurationController{
		configMapsAggregate: configMapsAggregate,
	}
}

type ReqModel struct {
	model.RequestBody
	Data map[string]string `json:"data,omitempty" protobuf:"bytes,1,opt,name=data"`
}

func (c *ConfigurationController) PostGitlab(req *ReqModel) (response model.Response, err error) {
	maps, err := c.configMapsAggregate.Create(constant.GitlabConstant, constant.TemplateDefaultNamespace, req.Data)
	response = new(model.BaseResponse)
	response.SetData(maps)
	return
}

func (c *ConfigurationController) PostDocker(req *ReqModel) (response model.Response, err error) {
	maps, err := c.configMapsAggregate.Create(constant.DockerConstant, constant.TemplateDefaultNamespace, req.Data)
	response = new(model.BaseResponse)
	response.SetData(maps)
	return
}
