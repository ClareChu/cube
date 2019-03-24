package aggregate

import (
	"github.com/magiconair/properties/assert"
	builder "hidevops.io/cube/console/pkg/builder/mocks"
	"hidevops.io/cube/console/pkg/command"
	"hidevops.io/cube/console/pkg/constant"
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/fake"
	"hidevops.io/cube/pkg/starter/cube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGatewayConfigCreate(t *testing.T) {
	clientSet := fake.NewSimpleClientset().CubeV1alpha1()
	gatewayClient := cube.NewGatewayConfig(clientSet)
	pb := new(builder.PipelineBuilder)
	gatewayAggregate := NewGatewayService(gatewayClient, pb)
	gc := &v1alpha1.GatewayConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       constant.GatewayConfigKind,
			APIVersion: constant.GatewayConfigApiVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "java",
			Namespace: constant.TemplateDefaultNamespace,
		},
	}
	_, err := gatewayClient.Create(gc)
	assert.Equal(t, nil, err)
	cmd := &command.GatewayConfig{}
	_, err = gatewayAggregate.Template(cmd)
	assert.Equal(t, nil, err)
	pb.On("Update", "", "demo", "createService", "fail", "").Return(nil)
	_, err = gatewayAggregate.Create("hello-world", "", "demo", "java", "v1", "dev")
	assert.Equal(t, nil, err)
}
