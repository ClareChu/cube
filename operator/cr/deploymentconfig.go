package cr

import (
	"hidevops.io/cube/operator/client"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentConfig struct {
	clientSet versioned.Interface
	Resource  string
}

const (
	DeploymentConfigResource = "pipelineconfigs"
)

func NewDeploymentConfig(clientSet versioned.Interface) CubeManagerInterface {
	return &DeploymentConfig{
		clientSet: clientSet,
		Resource:  DeploymentConfigResource,
	}
}

func (d *DeploymentConfig) create() {
	ide := &v1alpha1.DeploymentConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      IDEName,
			Namespace: Namespace,
		},
		Spec: v1alpha1.DeploymentConfigSpec{
			Container: corev1.Container{
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: 8443,
						Name:          "http-8443",
						Protocol:      corev1.ProtocolTCP,
					},
				},
			},
			EnvType: []string{
				"remoteDeploy",
				"deploy",
			},
			Profile: "dev",
			ReadinessProbe: corev1.Probe{
				Handler: corev1.Handler{
					Exec: &corev1.ExecAction{
						Command: []string{
							"curl",
							"--silent",
							"--show-error",
							"--fail",
							"http://localhost:8443",
						},
					},
				},
				FailureThreshold:    3,
				InitialDelaySeconds: 15,
				PeriodSeconds:       5,
				SuccessThreshold:    1,
				TimeoutSeconds:      1,
			},
		},
	}
	// create sonar pipeline config 流水线
	client.CreatDeploymentConfig(d.clientSet, ide)

	sonar := &v1alpha1.DeploymentConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      SonarName,
			Namespace: Namespace,
		},
		Spec: v1alpha1.DeploymentConfigSpec{
			Container: corev1.Container{
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: 8080,
						Name:          "http-8080",
						Protocol:      corev1.ProtocolTCP,
					},
				},
			},
			EnvType: []string{
				"remoteDeploy",
				"deploy",
			},
			Profile: "dev",
			ReadinessProbe: corev1.Probe{
				Handler: corev1.Handler{
					Exec: &corev1.ExecAction{
						Command: []string{
							"curl",
							"--silent",
							"--show-error",
							"--fail",
							"http://localhost:8080/actuator/health",
						},
					},
				},
				FailureThreshold:    3,
				InitialDelaySeconds: 5,
				PeriodSeconds:       5,
				SuccessThreshold:    1,
				TimeoutSeconds:      1,
			},
		},
	}
	// create sonar pipeline config 流水线
	client.CreatDeploymentConfig(d.clientSet, sonar)
}
