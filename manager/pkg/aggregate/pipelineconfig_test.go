package aggregate

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/cube/manager/pkg/aggregate/mocks"
	"hidevops.io/cube/manager/pkg/command"
	"hidevops.io/cube/manager/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestPipelineConfigTemplate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	pipelineConfig := cube.NewPipelineConfig(clientSet)
	pipelineAggregate := new(mocks.PipelineAggregate)
	pa := NewPipelineConfigService(pipelineConfig, pipelineAggregate)
	cdc := &command.PipelineConfigTemplate{}
	_, err := pa.NewPipelineConfigTemplate(cdc)
	assert.Equal(t, nil, err)
}

func TestPipelineConfigCreate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	pipelineConfig := cube.NewPipelineConfig(clientSet)
	pipelineAggregate := new(mocks.PipelineAggregate)
	pa := NewPipelineConfigService(pipelineConfig, pipelineAggregate)
	dc := &v1alpha1.PipelineConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "java",
			Namespace: constant.TemplateDefaultNamespace,
		},
	}
	_, err := pipelineConfig.Create(dc)
	_, err = pa.Create("hello-world", "hello-world-1", dc)
	assert.Equal(t, nil, err)

	_, err = pa.Get("hello-world", "hello-world-1")
	assert.Equal(t, nil, err)
}
