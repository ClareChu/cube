package aggregate

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/mio/console/pkg/aggregate/mocks"
	"hidevops.io/mio/console/pkg/command"
	"hidevops.io/mio/console/pkg/constant"
	"hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"hidevops.io/mio/pkg/client/clientset/versioned/fake"
	"hidevops.io/mio/pkg/starter/mio"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestPipelineConfigTemplate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	pipelineConfig := mio.NewPipelineConfig(clientSet)
	pipelineAggregate := new(mocks.PipelineAggregate)
	pa := NewPipelineConfigService(pipelineConfig, pipelineAggregate)
	cdc := &command.PipelineConfigTemplate{}
	_, err := pa.NewPipelineConfigTemplate(cdc)
	assert.Equal(t, nil, err)
}

func TestPipelineConfigCreate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().MioV1alpha1()
	pipelineConfig := mio.NewPipelineConfig(clientSet)
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
