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
)

type ImagePushInterface interface {
	Handle(build *v1alpha1.Build) error
}

type ImagePush struct {
	ImagePushInterface
	buildConfigService service.BuildConfigService
	buildPackInterface dispatch.BuildPackInterface
}

func init() {
	app.Register(NewImagePushService)
}

const PushImage = "pushImage"

func NewImagePushService(buildConfigService service.BuildConfigService,
	buildPackInterface dispatch.BuildPackInterface) CodeInterface {
	imagePush := &ImagePush{
		buildConfigService: buildConfigService,
		buildPackInterface: buildPackInterface,
	}
	buildPackInterface.Add(PushImage, imagePush)
	return imagePush
}

func (i *ImagePush) Handle(build *v1alpha1.Build) error {
	cmd := &command.ImagePushCommand{
		Tags:      []string{build.Spec.Tags[0] + ":" + build.ObjectMeta.Labels[constant.Number], build.Spec.Tags[0] + ":" + constant.DefaultTag},
		Name:      build.Name,
		Namespace: build.Namespace,
		Username:  build.Spec.DockerAuthConfig.Username,
		Password:  build.Spec.DockerAuthConfig.Password,
		ImageName: build.ObjectMeta.Labels["name"],
	}
	log.Infof("ImagePush :%v", cmd)
	err := i.buildConfigService.ImagePush(fmt.Sprintf("%s.%s.svc", build.Name, build.Namespace), constant.DefaultPort, cmd)
	return err
}
