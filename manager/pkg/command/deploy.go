package command

import (
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/extensions/v1beta1"
)

type DeployData struct {
	Name          string
	Namespace     string
	Labels        map[string]string
	NodeSelector  map[string]string
	Container     corev1.Container
	Volumes       []corev1.Volume
	Replicas      *int32
	Version       string
	Strategy      appsv1.DeploymentStrategy
	InitContainer corev1.Container
	ForceUpdate   bool
}
