package aggregate

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/hioak/starter/kube"
	"hidevops.io/mio/console/pkg/aggregate/mocks"
	builder "hidevops.io/mio/console/pkg/builder/mocks"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/client/clientset/versioned/fake"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeFake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestDeploymentCreate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	deployment := mio.NewDeployment(clientSet)
	remoteAggregate := new(mocks.RemoteAggregate)
	client := kubeFake.NewSimpleClientset()
	deploymentClient := kube.NewDeployment(client)
	deployBuilder := new(builder.DeploymentBuilder)
	pipelineBuilder := new(builder.PipelineBuilder)
	buildConfigAggregate := NewDeploymentService(deployment, remoteAggregate, deploymentClient, deployBuilder, pipelineBuilder)
	dc := &v1alpha1.DeploymentConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	_, err := buildConfigAggregate.Create(dc, "hello-world", "v1", "2")
	assert.Equal(t, nil, err)
}

func TestDeploymentCreateApp(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	deployment := mio.NewDeployment(clientSet)
	remoteAggregate := new(mocks.RemoteAggregate)
	client := kubeFake.NewSimpleClientset()
	deploymentClient := kube.NewDeployment(client)
	deployBuilder := new(builder.DeploymentBuilder)
	pipelineBuilder := new(builder.PipelineBuilder)
	buildConfigAggregate := NewDeploymentService(deployment, remoteAggregate, deploymentClient, deployBuilder, pipelineBuilder)
	d := &v1alpha1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hello-world",
			Namespace: "demo",
		},
	}
	deployBuilder.On("Update", "hello-world", "demo", "deploy", "success").Return(nil)
	err := buildConfigAggregate.CreateApp(d)
	assert.Equal(t, nil, err)
}

func TestDeploymentSelector(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	deployment := mio.NewDeployment(clientSet)
	remoteAggregate := new(mocks.RemoteAggregate)
	client := kubeFake.NewSimpleClientset()
	deploymentClient := kube.NewDeployment(client)
	deployBuilder := new(builder.DeploymentBuilder)
	pipelineBuilder := new(builder.PipelineBuilder)
	buildConfigAggregate := NewDeploymentService(deployment, remoteAggregate, deploymentClient, deployBuilder, pipelineBuilder)
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
