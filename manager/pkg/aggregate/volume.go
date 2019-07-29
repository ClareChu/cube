package aggregate

import (
	"hidevops.io/cube/manager/pkg/builder"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VolumeAggregate interface {
	Create(params *command.PipelineReqParams) error
	CreateVolume(namespace string, volume v1alpha1.Volumes) error
}

type volumeServiceImpl struct {
	persistentVolumeClaim *kube.PersistentVolumeClaim
	pipelineBuilder       builder.PipelineBuilder
}

func init() {
	app.Register(NewVolumeService)
}

func NewVolumeService(persistentVolumeClaim *kube.PersistentVolumeClaim, pipelineBuilder builder.PipelineBuilder) VolumeAggregate {
	return &volumeServiceImpl{
		persistentVolumeClaim: persistentVolumeClaim,
		pipelineBuilder:       pipelineBuilder,
	}
}

func (v *volumeServiceImpl) Create(params *command.PipelineReqParams) (err error) {
	if params.Volumes.Name == "" {
		err = v.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.CreateService, constant.Success, "")
		return err
	}
	err = v.CreateVolume(params.Namespace, params.Volumes)
	if err != nil {
		v.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.CreateService, constant.Fail, "")
		return err
	}
	err = v.pipelineBuilder.Update(params.PipelineName, params.Namespace, constant.CreateService, constant.Success, "")
	return err
}

type DataSize string

const (
	MByte    DataSize = "Mi"
	GigaByte DataSize = "Gi"
)

func (v *volumeServiceImpl) CreateVolume(namespace string, volume v1alpha1.Volumes) error {
	options := v1.GetOptions{}
	_, err := v.persistentVolumeClaim.Get(volume.Name, namespace, options)
	if err == nil {
		log.Debugf("pvc already existed !!!")
		return err
	}
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      volume.Name,
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
			StorageClassName: &volume.StorageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(volume.Size),
				},
			},
		},
	}
	_, err = v.persistentVolumeClaim.Create(pvc)
	return err
}

func Validate(volumes v1alpha1.Volumes) bool {
	return false
}
