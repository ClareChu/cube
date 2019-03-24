package aggregate

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/cube/console/pkg/aggregate/mocks"
	builder "hidevops.io/cube/console/pkg/builder/mocks"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/cube/console/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"hidevops.io/hioak/starter/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeFake "k8s.io/client-go/kubernetes/fake"
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
	_, err = buildConfigAggregate.Create("hello-world", "hello-world-1", "demo", "java", "v1", "1", "dev")
	assert.Equal(t, nil, err)
}
