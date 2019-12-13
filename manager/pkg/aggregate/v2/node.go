package v2

import (
	"hidevops.io/cube/manager/pkg/aggregate/dispatch"
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
)

type NodeInterface interface {
	Handle(build *v1alpha1.Build) error
}

type Node struct {
	NodeInterface
	buildNode          builder.BuildNode
	build              BuildInterface
	buildPackInterface dispatch.BuildPackInterface
}

const DeployNode = "deployNode"

func init() {
	app.Register(NewNodeService)
}

func NewNodeService(buildNode builder.BuildNode,
	build BuildInterface,
	buildPackInterface dispatch.BuildPackInterface) NodeInterface {
	node := &Node{
		buildNode:          buildNode,
		build:              build,
		buildPackInterface: buildPackInterface,
	}
	buildPackInterface.Add(DeployNode, node)
	return node
}

func (n *Node) Handle(build *v1alpha1.Build) error {
	phase := constant.Success
	volumes, volumeMounts := volume(build)
	command := &command.DeployNode{
		DeployData: kube.DeployData{
			Name:      build.Name,
			NameSpace: build.Namespace,
			Replicas:  build.Spec.DeployData.Replicas,
			Labels: map[string]string{
				constant.App:  build.Name,
				constant.Name: build.ObjectMeta.Labels[constant.Name],
			},
			Image:        build.Spec.BaseImage,
			Ports:        build.Spec.DeployData.Ports,
			Envs:         build.Spec.DeployData.Envs,
			NodeName:     build.Spec.DeployData.Envs["NODE_NAME"],
			Volumes:      volumes,
			VolumeMounts: volumeMounts,
		},
	}
	_, err := n.buildNode.Start(command)
	if err != nil {
		log.Errorf("deploy agent err :%v", err)
		return err
	}
	err = n.build.WatchPod(build.Name, build.Namespace)
	if err != nil {
		phase = constant.Fail
	}

	err = n.build.Update(build, constant.DeployNode, phase)
	return err
}

func volume(build *v1alpha1.Build) (volumes []corev1.Volume, volumeMounts []corev1.VolumeMount) {
	for _, hostPathVolume := range build.Spec.DeployData.HostPathVolumes {
		volume := corev1.Volume{
			Name: hostPathVolume.Name,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: hostPathVolume.VolumeSource,
				},
			},
		}
		volumeMount := corev1.VolumeMount{
			Name:      hostPathVolume.Name,
			ReadOnly:  hostPathVolume.ReadOnly,
			MountPath: hostPathVolume.MountPath,
			SubPath:   hostPathVolume.SubPath,
		}
		volumes = append(volumes, volume)
		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumes, volumeMounts
}
