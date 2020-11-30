package cr

import (
	"hidevops.io/cube/operator/client"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceConfig struct {
	clientSet versioned.Interface
	Resource  string
}

const (
	ServiceConfigResource = "serviceconfigs"
)

func NewServiceConfig(clientSet versioned.Interface) CubeManagerInterface {
	return &ServiceConfig{
		clientSet: clientSet,
		Resource:  ServiceConfigResource,
	}
}

func (s *ServiceConfig) create() {
	ide := &v1alpha1.ServiceConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      IDEName,
			Namespace: Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 8443,
					TargetPort: intstr.IntOrString{
						IntVal: 8443,
					},
					Name:     "http-8443",
					Protocol: corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	// create ide pipeline config 流水线
	//client.CreateRestClient(s.clientSet, ide, s.Resource)
	client.CreatServiceConfig(s.clientSet, ide)
	sonar := &v1alpha1.ServiceConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      SonarName,
			Namespace: Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 8080,
					TargetPort: intstr.IntOrString{
						IntVal: 8080,
					},
					Name:     "http-8080",
					Protocol: corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	// create sonar pipeline config 流水线
	//client.CreateRestClient(s.clientSet, sonar, s.Resource)
	client.CreatServiceConfig(s.clientSet, sonar)
}
