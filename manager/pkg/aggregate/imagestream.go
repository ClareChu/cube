package aggregate

import (
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
)

type ImageStreamAggregate interface {
	CreateImage(params *command.PipelineReqParams) error
	Create(name, namespace string, images string) error
}

type imageStreamServiceImpl struct {
	imageStream     *cube.ImageStream
	pipelineBuilder builder.PipelineBuilder
}

func init() {
	app.Register(NewImageStreamService)
}

func NewImageStreamService(imageStream *cube.ImageStream, pipelineBuilder builder.PipelineBuilder) ImageStreamAggregate {
	return &imageStreamServiceImpl{
		imageStream:     imageStream,
		pipelineBuilder: pipelineBuilder,
	}
}

const (
	Latest     = "latest"
	Generation = "1"
)

func (i *imageStreamServiceImpl) CreateImage(params *command.PipelineReqParams) error {
	err := i.Create(params.Name, params.Namespace, params.Container.Image)
	if err != nil {
		i.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.ImageStream, constant.Fail, "")
		return err
	}
	err = i.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.ImageStream, constant.Success, "")
	return err
}

func (i *imageStreamServiceImpl) Create(name, namespace string, images string) error {
	tag := strings.Split(images, ":")[1]
	t := time.Now()
	image, err := i.imageStream.Get(name, namespace)
	stream := &v1alpha1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
	}
	spec := v1alpha1.ImageStreamSpec{
		DockerImageRepository: images,
		Tags: map[string]v1alpha1.Tag{
			tag: v1alpha1.Tag{
				Created:              t.UTC().Format(time.UnixDate),
				DockerImageReference: images,
				Generation:           "1",
				Image:                strings.Split(images, ":")[1],
			},
			Latest: v1alpha1.Tag{
				Created:              t.UTC().Format(time.UnixDate),
				DockerImageReference: images,
				Generation:           "1",
				Image:                strings.Split(images, ":")[1],
			},
		},
	}
	stream.Spec = spec
	if err != nil {
		log.Debugf("get image: %v", err)
		_, err := i.imageStream.Create(stream)
		log.Infof("create image stream error :%v", err)
		return err
	}
	delete(image.Spec.Tags, tag)
	delete(image.Spec.Tags, Latest)
	image.Spec.DockerImageRepository = images
	image.Spec.Tags[tag] = v1alpha1.Tag{
		Created:              t.UTC().Format(time.UnixDate),
		DockerImageReference: images,
		Generation:           "1",
		Image:                strings.Split(images, ":")[1],
	}
	image.Spec.Tags[Latest] = v1alpha1.Tag{
		Created:              t.UTC().Format(time.UnixDate),
		DockerImageReference: images,
		Generation:           "1",
		Image:                strings.Split(images, ":")[1],
	}
	_, err = i.imageStream.Update(name, namespace, image)
	return err
}
