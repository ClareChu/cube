package command

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type DeployData struct {
	Name           string
	Namespace      string
	Labels         map[string]string
	NodeSelector   map[string]string
	Container      corev1.Container
	Volumes        []corev1.Volume
	Replicas       *int32
	Version        string
	Strategy       appsv1.DeploymentStrategy
	InitContainer  corev1.Container
	ForceUpdate    bool
	ReadinessProbe corev1.Probe
	LivenessProbe  corev1.Probe
}
