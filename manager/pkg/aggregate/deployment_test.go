package aggregate

import (
	"errors"
	"github.com/magiconair/properties/assert"
	"hidevops.io/cube/manager/pkg/aggregate/mocks"
	builder "hidevops.io/cube/manager/pkg/builder/mocks"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
)

func TestDeploymentCreate(t *testing.T) {
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	deployment := cube.NewDeployment(clientSet)
	remoteAggregate := new(mocks.RemoteAggregate)
	deployBuilder := new(builder.DeploymentBuilder)
	pipelineBuilder := new(builder.PipelineBuilder)
	deploy := new(builder.DeploymentConfigBuilder)
	tag := new(mocks.TagAggregate)
	buildConfigAggregate := NewDeploymentService(deployment, remoteAggregate, deployBuilder, pipelineBuilder, deploy, tag)
	dc := &v1alpha1.DeploymentConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	_, err := buildConfigAggregate.Create(dc, "hello-world", "v1", "2")
	assert.Equal(t, errors.New("pod query timeout 10 minutes"), err)
}

func TestDeploymentSelector(t *testing.T) {
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	deployment := cube.NewDeployment(clientSet)
	remoteAggregate := new(mocks.RemoteAggregate)
	deployBuilder := new(builder.DeploymentBuilder)
	pipelineBuilder := new(builder.PipelineBuilder)
	deploy := new(builder.DeploymentConfigBuilder)
	tag := new(mocks.TagAggregate)
	buildConfigAggregate := NewDeploymentService(deployment, remoteAggregate, deployBuilder, pipelineBuilder, deploy, tag)
	d := &v1alpha1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
		Status: v1alpha1.DeploymentStatus{
			Stages: []v1alpha1.Stages{v1alpha1.Stages{
				Name: constant.Deploy,
			}},
		},
	}
	deployBuilder.On("Update", "hello-world", "demo", "deploy", "success").Return(nil)
	err := buildConfigAggregate.Selector(d)
	assert.Equal(t, nil, err)
}

func TestDeployCreate(t *testing.T) {
	deployBuilder := new(builder.DeploymentBuilder)
	deploy := new(builder.DeploymentConfigBuilder)
	tag := new(mocks.TagAggregate)
	buildConfigAggregate := NewDeploymentService(nil, nil, deployBuilder, nil, deploy, tag)
	d := &v1alpha1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
		Status: v1alpha1.DeploymentStatus{
			Stages: []v1alpha1.Stages{v1alpha1.Stages{
				Name: constant.Deploy,
			}},
		},
	}
	deployBuilder.On("CreateApp", d).Return(nil)
	err := buildConfigAggregate.CreateDeployment(d)
	assert.Equal(t, nil, err)
	err = os.Setenv("IS_OPENSHIFT", "1")
	assert.Equal(t, nil, err)
	deploy.On("CreateApp", d).Return(nil)
	err = buildConfigAggregate.CreateDeployment(d)
	assert.Equal(t, nil, err)
}
