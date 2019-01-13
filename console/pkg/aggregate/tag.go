package aggregate

import (
	"fmt"
	"github.com/prometheus/common/log"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/mio/console/pkg/builder"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	miov1alpha1 "hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/starter/mio"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type TagAggregate interface {
	TagImage(deploy *v1alpha1.Deployment) error
}

func init() {
	app.Register(NewTagService)
}

type Tag struct {
	TagAggregate
	imageStream       *mio.ImageStream
	deploymentBuilder builder.DeploymentBuilder
}

func NewTagService(imageStream *mio.ImageStream, deploymentBuilder builder.DeploymentBuilder) TagAggregate {
	return &Tag{
		imageStream:       imageStream,
		deploymentBuilder: deploymentBuilder,
	}
}

func (t *Tag) TagImage(deploy *v1alpha1.Deployment) error {
	i, err := t.imageStream.Get(deploy.Labels[constant.DeploymentConfig], deploy.Namespace)
	if err != nil {
		log.Errorf("get image: %v", err)
		return err
	}
	tag := i.Spec.Tags[constant.Latest]
	n := GetNamespace(deploy.Namespace, deploy.Spec.Profile)
	is, err := t.imageStream.Get(deploy.Labels[constant.DeploymentConfig], n)
	if err == nil {
		imageStream := &v1alpha1.ImageStream{
			ObjectMeta: is.ObjectMeta,
			Spec: v1alpha1.ImageStreamSpec{
				Tags: map[string]v1alpha1.Tag{
					constant.Latest: tag,
				},
				DockerImageRepository: strings.Split(i.Spec.DockerImageRepository, ":")[0] + ":latest" ,
			},
		}
		_, err = t.imageStream.Update(deploy.Labels[constant.DeploymentConfig], n, imageStream)
		err = t.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.RemoteDeploy, constant.Success)
		return err
	}
	stream := &miov1alpha1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploy.Labels[constant.DeploymentConfig],
			Namespace: n,
			Labels: map[string]string{
				"app": deploy.Labels[constant.DeploymentConfig],
			},
		},
		Spec: miov1alpha1.ImageStreamSpec{
			DockerImageRepository: strings.Split(i.Spec.DockerImageRepository, ":")[0] + ":latest" ,
			Tags: map[string]miov1alpha1.Tag{
				constant.Latest: tag,
			},
		},
	}
	is, err = t.imageStream.Create(stream)
	err = t.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.RemoteDeploy, constant.Success)
	return err
}

func GetNamespace(space, profile string) string {
	if profile == "" {
		return space
	}
	return fmt.Sprintf("%s-%s", space, profile)
}
