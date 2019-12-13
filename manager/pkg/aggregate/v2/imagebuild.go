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
	"hidevops.io/hiboot/pkg/utils/idgen"
)

type ImageBuildInterface interface {
	Handle(build *v1alpha1.Build) error
}

type ImageBuild struct {
	ImageBuildInterface
	buildConfigService service.BuildConfigService
	buildPackInterface dispatch.BuildPackInterface
}

func init() {
	app.Register(NewImageBuildService)
}

const BuildImage = "buildImage"

func NewImageBuildService(buildConfigService service.BuildConfigService,
	buildPackInterface dispatch.BuildPackInterface) CodeInterface {
	imageBuild := &ImageBuild{
		buildConfigService: buildConfigService,
		buildPackInterface: buildPackInterface,
	}
	buildPackInterface.Add(BuildImage, imageBuild)
	return imageBuild
}

func (i *ImageBuild) Handle(build *v1alpha1.Build) error {
	id, err := idgen.NextString()
	if err != nil {
		log.Errorf("id err :{}", id)
	}
	cmd := &command.ImageBuildCommand{
		App:        build.Spec.App,
		S2IImage:   build.Spec.BaseImage,
		Tags:       []string{build.Spec.Tags[0] + ":" + build.ObjectMeta.Labels[constant.Number], build.Spec.Tags[0] + ":" + constant.DefaultTag},
		DockerFile: build.Spec.DockerFile,
		Name:       build.Name,
		Namespace:  build.Namespace,
		Username:   build.Spec.DockerAuthConfig.Username,
		Password:   build.Spec.DockerAuthConfig.Password,
	}
	log.Infof("build ImageBuildBuild :%v", cmd)
	err = i.buildConfigService.ImageBuild(fmt.Sprintf("%s.%s.svc", build.Name, build.Namespace), constant.DefaultPort, cmd)
	return err
}
