package aggregate

import (
	"fmt"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/docker"
	"hidevops.io/hioak/starter/kube"
	"os"
)

const ImageUsername = "unused"

type RemoteAggregate interface {
	TagImage(deploy *v1alpha1.Deployment) error
}

func init() {
	app.Register(NewRemoteService)
}

type Remote struct {
	RemoteAggregate
	imageClient       *docker.ImageClient
	token             kube.Token
	deploymentBuilder builder.DeploymentBuilder
}

func NewRemoteService(imageClient *docker.ImageClient, token kube.Token, deploymentBuilder builder.DeploymentBuilder) RemoteAggregate {
	return &Remote{
		imageClient:       imageClient,
		token:             token,
		deploymentBuilder: deploymentBuilder,
	}
}

func (r *Remote) TagImage(deploy *v1alpha1.Deployment) error {
	log.Infof("===================deploy :%v", deploy)
	token := fmt.Sprintf("%v", r.token)
	username := ""
	password := ""
	if deploy.Spec.DockerAuthConfig.Username == ImageUsername || deploy.Spec.DockerAuthConfig.Username == "" {
		username = ImageUsername
		password = token
	} else {
		username = deploy.Spec.DockerAuthConfig.Username
		password = deploy.Spec.DockerAuthConfig.Password
	}
	//TODO pull image
	image := &docker.Image{
		Username:  username,
		Password:  password,
		FromImage: deploy.Spec.FromRegistry + "/" + deploy.Namespace + "/" + deploy.ObjectMeta.Labels[constant.DeploymentConfig],
		Tag:       deploy.ObjectMeta.Labels[constant.BuildVersion],
	}
	log.Debugf(" get tag image :%v", image)
	err := r.imageClient.PullImage(image)
	if err != nil {
		log.Info("pull image error :", err)
		err = r.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.RemoteDeploy, constant.Fail)
		return err
	}
	//TODO get image
	s, err := r.imageClient.GetImage(image)
	if err != nil {
		log.Info("get image error :", err)
		err = r.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.RemoteDeploy, constant.Fail)
		return err
	}
	image.FromImage = deploy.Spec.FromRegistry + "/" + deploy.Namespace + "-" + deploy.Spec.Profile + "/" + deploy.ObjectMeta.Labels[constant.DeploymentConfig]
	//TODO tag IMAGE
	err = r.imageClient.TagImage(image, s.ID)
	if err != nil {
		log.Info("tag image error :", err)
		err = r.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.RemoteDeploy, constant.Fail)
		return err
	}
	//TODO PUSH IMAGE
	err = r.imageClient.PushImage(image)
	if err != nil {
		log.Info("tag image error :", err)
		err = r.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.RemoteDeploy, constant.Fail)
		return err
	}
	err = r.deploymentBuilder.Update(deploy.Name, deploy.Namespace, constant.RemoteDeploy, constant.Success)
	return err
}

//GetRemoteTag get tag if
func GetRemoteTag(version string) (tag string) {
	tag = os.Getenv("Image_Tag")
	if tag == "" {
		tag = version
		return
	}
	return
}
