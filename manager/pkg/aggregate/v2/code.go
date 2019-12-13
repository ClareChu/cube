package v2

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/aggregate/dispatch"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/manager/pkg/service"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
)

type CodeInterface interface {
	Handle(build *v1alpha1.Build) error
}

type Code struct {
	CodeInterface
	configMaps         *kube.ConfigMaps
	buildConfigService service.BuildConfigService
	buildPackInterface dispatch.BuildPackInterface
}

const CLONE = "clone"

func init() {
	app.Register(NewCodeService)
}

func NewCodeService(configMaps *kube.ConfigMaps,
	buildConfigService service.BuildConfigService,
	buildPackInterface dispatch.BuildPackInterface) CodeInterface {
	code := &Code{
		configMaps:         configMaps,
		buildConfigService: buildConfigService,
		buildPackInterface: buildPackInterface,
	}
	buildPackInterface.Add(CLONE, code)
	return code
}

func (c *Code) Handle(build *v1alpha1.Build) error {
	configMaps, err := c.configMaps.Get(constant.GitlabConstant, constant.TemplateDefaultNamespace)
	if err != nil {
		log.Errorf("build get configMaps: %v", err)
		return err
	}
	baseUrl := build.Spec.CloneConfig.Url
	if baseUrl == "" {
		baseUrl = configMaps.Data[constant.BaseUrl]
	}
	log.Debugf("baseUrl: %v", baseUrl)
	command := &command.SourceCodePullCommand{
		CloneType: build.Spec.CloneType,
		Url:       fmt.Sprintf("%s/%s/%s.git", baseUrl, build.Spec.Project, build.Spec.AppRoot),
		Branch:    build.Spec.CloneConfig.Branch,
		DstDir:    build.Spec.CloneConfig.DstDir,
		Username:  build.Spec.CloneConfig.Username,
		Password:  build.Spec.CloneConfig.Password,
		Namespace: build.Namespace,
		Name:      build.Name,
	}

	err = c.buildConfigService.SourceCodePull(fmt.Sprintf("%s.%s.svc", build.Name, build.Namespace), constant.DefaultPort, command)
	return err
}
