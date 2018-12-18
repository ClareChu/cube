package aggregate

import (
	"errors"
	"github.com/magiconair/properties/assert"
	"hidevops.io/mio/console/pkg/aggregate/mocks"
	builder "hidevops.io/mio/console/pkg/builder/mocks"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/client/clientset/versioned/fake"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
)

func TestDeploymentCreate(t *testing.T) {
	os.Setenv("KUBE_WATCH_TIMEOUT", "1")
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	deployment := mio.NewDeployment(clientSet)
	remoteAggregate := new(mocks.RemoteAggregate)
	deployBuilder := new(builder.DeploymentBuilder)
	pipelineBuilder := new(builder.PipelineBuilder)
	deploy := new(builder.DeploymentConfigBuilder)
	buildConfigAggregate := NewDeploymentService(deployment, remoteAggregate, deployBuilder, pipelineBuilder, deploy)
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
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	deployment := mio.NewDeployment(clientSet)
	remoteAggregate := new(mocks.RemoteAggregate)
	deployBuilder := new(builder.DeploymentBuilder)
	pipelineBuilder := new(builder.PipelineBuilder)
	deploy := new(builder.DeploymentConfigBuilder)
	buildConfigAggregate := NewDeploymentService(deployment, remoteAggregate, deployBuilder, pipelineBuilder, deploy)
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
