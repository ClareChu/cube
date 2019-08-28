package aggregate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"hidevops.io/cube/manager/pkg/aggregate/mocks"
	builder "hidevops.io/cube/manager/pkg/builder/mocks"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hioak/starter/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeFake "k8s.io/client-go/kubernetes/fake"
	"sort"
	"testing"
)

func TestDeploymentConfigTemplate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	deploymentConfig := cube.NewDeploymentConfig(clientSet)
	deploymentAggregate := new(mocks.DeploymentAggregate)
	pipelineBuilder := new(builder.PipelineBuilder)
	client := kubeFake.NewSimpleClientset()
	deploymentClient := kube.NewDeployment(client)
	buildConfigAggregate := NewDeploymentConfigService(deploymentConfig, deploymentClient, pipelineBuilder, deploymentAggregate)
	cdc := &command.DeploymentConfig{}
	_, err := buildConfigAggregate.Template(cdc)
	assert.Equal(t, nil, err)
}

func TestDeploymentConfigCreate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	deploymentConfig := cube.NewDeploymentConfig(clientSet)
	deploymentAggregate := new(mocks.DeploymentAggregate)
	pipelineBuilder := new(builder.PipelineBuilder)
	client := kubeFake.NewSimpleClientset()
	deploymentClient := kube.NewDeployment(client)
	buildConfigAggregate := NewDeploymentConfigService(deploymentConfig, deploymentClient, pipelineBuilder, deploymentAggregate)
	dc := &v1alpha1.DeploymentConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "java",
			Namespace: constant.TemplateDefaultNamespace,
		},
	}
	_, err := deploymentConfig.Create(dc)
	d := &v1alpha1.DeploymentConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       "DeploymentConfig.cube.io/v1alpha1",
			APIVersion: "a1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
		Spec: v1alpha1.DeploymentConfigSpec{
			Profile: "dev",
		},
		Status: v1alpha1.DeploymentConfigStatus{
			LastVersion: 1,
		},
	}
	deploymentAggregate.On("Create", d, "hello-world-1", "v1", "1").Return(nil, nil)
	param := &command.PipelineReqParams{
		Name:         "hello-world",
		PipelineName: "hello-world-1",
		Namespace:    "demo",
		EventType:    "java",
		Version:      "v1",
		Profile:      "dev",
	}
	_, err = buildConfigAggregate.Create(param, "1")
	assert.Equal(t, nil, err)
}

func TestDeploymentConfig(t *testing.T) {
	buildConfigAggregate := NewDeploymentConfigService(nil, nil, nil, nil)
	deploy := &v1alpha1.DeploymentConfig{
		Spec: v1alpha1.DeploymentConfigSpec{
			Container: corev1.Container{
				Env: []corev1.EnvVar{
					corev1.EnvVar{
						Name: "c",
						Value: "d",
					},
				},
				Command: []string{"a", "s"},
				Ports: []corev1.ContainerPort{
					corev1.ContainerPort{
						Name: "http-8082",
						Protocol: corev1.ProtocolTCP,
						ContainerPort: 8082,
					},
				},
			},
		},
	}
	template := &v1alpha1.DeploymentConfig{
		Spec: v1alpha1.DeploymentConfigSpec{
			Container: corev1.Container{
				Env: []corev1.EnvVar{
					corev1.EnvVar{
						Name: "c",
						Value: "d",
					},
				},
				Command: []string{"a", "s"},
				Ports: []corev1.ContainerPort{
					corev1.ContainerPort{
						Name: "http-8082",
						Protocol: corev1.ProtocolTCP,
						ContainerPort: 8082,
					},
				},
			},
		},
	}
	param := &command.PipelineReqParams{
		Container: corev1.Container{
			Command: []string{"a", "b"},
			Ports: []corev1.ContainerPort{
				corev1.ContainerPort{
					Name: "http-8081",
					Protocol: corev1.ProtocolTCP,
					ContainerPort: 8081,
				},
			},
		},
	}
	buildConfigAggregate.InitDeployConfig(deploy, template, param)
	assert.Equal(t, deploy.Spec.Container.Command[1], param.Container.Command[1])
	assert.Equal(t, deploy.Spec.Container.Command[0], param.Container.Command[0])
	assert.Equal(t, deploy.Spec.Container.Env[0].Name, template.Spec.Container.Env[0].Name)
	assert.Equal(t, deploy.Spec.Container.Env[0].Value, template.Spec.Container.Env[0].Value)
}

func TestSort(t *testing.T)  {
	s := []string{"a","c","s","d","w","b", "ac", "ab"}
	sort.Strings(s)
	fmt.Printf("%v", s)
}