package service

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	extensionsV1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type DeployService interface {
	Create(dd *DeployData) (*extensionsV1beta1.Deployment, error)
}
type Deploy struct {
	DeployService
	deployment *kube.Deployment
}

func init() {
	app.Register(NewDeploy)
}

func NewDeploy(deployment *kube.Deployment) DeployService {
	return &Deploy{
		deployment: deployment,
	}
}

type DeployData struct {
	Container corev1.Container
}

func (d *Deploy) Create(dd *DeployData) (*extensionsV1beta1.Deployment, error) {
	dpm := &extensionsV1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      dd.Name,
			Namespace: dd.NameSpace,
		},
		Spec: extensionsV1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Strategy: extensionsV1beta1.DeploymentStrategy{
				Type: extensionsV1beta1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &extensionsV1beta1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(1),
					},
				},
			},
			RevisionHistoryLimit: int32Ptr(10),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   dd.Name,
					Labels: dd.Labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: dd.NodeSelector,
					Containers: []corev1.Container{
						dd.Container,
					},
				},
			},
		},
	}

	d.deployment.Create(dpm)

}
